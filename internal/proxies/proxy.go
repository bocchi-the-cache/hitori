package proxies

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"

	"github.com/bocchi-the-cache/hitori/pkg/cache"
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/bocchi-the-cache/hitori/pkg/logger"
	"github.com/bocchi-the-cache/hitori/pkg/origin"
)

const ServerToken = "hitori-cache-server"

var DefaultProxy *Proxy

func Init(config *config.Config) {
	DefaultProxy = NewProxy(config.Server.Port, cache.DefaultCache, origin.DefaultOrigin)
	logger.Info("proxy init successfully", DefaultProxy.addr)
}

func Serve() error {
	return DefaultProxy.Serve()
}

type Proxy struct {
	addr string
	s    *fasthttp.Server
	c    cache.Cache
	o    *origin.Origin
}

// NewProxy
// TODO: use buildOption to support complex proxies settings
func NewProxy(port int, ca cache.Cache, ori *origin.Origin) *Proxy {
	p := new(Proxy)
	p.addr = fmt.Sprintf(":%d", port)
	p.s = &fasthttp.Server{
		Handler: p.ProxyHandler,
		Name:    ServerToken,
	}
	p.o = ori
	p.c = ca
	return p
}

func (p *Proxy) ListenAddr() string { return p.addr }

func (p *Proxy) Serve() error {
	logger.Info("proxy server start", p.addr)
	return p.s.ListenAndServe(p.addr)
}

func (p *Proxy) ProxyHandler(ctx *fasthttp.RequestCtx) {
	logger.Debugf("receive client request uri: %v", ctx.Request.URI())
	ctx.Response.Header.Add("X-Proxy", ServerToken)

	// TODO: map use buffer instead of []byte
	// TODO: cache key management

	err := p.ServeByCache(ctx)
	if err != nil {
		logger.Errorf("cache serve error: %v", err)
	}

}

func (p *Proxy) ServeByCache(ctx *fasthttp.RequestCtx) error {
	headerKey := ctx.Request.URI().String() + "_header"
	HeaderData, err := p.c.Get(headerKey)

	// cache MISS or error
	if err != nil {
		logger.Errorf("cache header get error, url: %v, err: %v", ctx.Request.URI(), err)
		return p.ServeByOrigin(ctx)
	}
	if HeaderData == nil {
		logger.Debugf("cache miss, url: %v", ctx.Request.URI())
		return p.ServeByOrigin(ctx)
	}
	logger.Debugf("read header data from cache success, length:%v, key: %v", len(HeaderData), headerKey)

	// cache HIT
	bodyKey := ctx.Request.URI().String() + "_body"
	BodyData, err := p.c.Get(bodyKey)
	if err != nil {
		logger.Errorf("cache body get error, url: %v, err: %v", ctx.Request.URI(), err)
		return p.ServeByOrigin(ctx)
	}
	if BodyData == nil {
		logger.Errorf("cache header hit but body missed, url: %v", ctx.Request.URI())
		return p.ServeByOrigin(ctx)
	}
	logger.Debugf("read body data from cache success, length:%v, key: %v", len(BodyData), bodyKey)

	// parse header from HeaderData
	// TODO: unfortunately, fasthttp does not support to set header from []byte
	normalizedResp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(HeaderData)), nil)
	for k, v := range normalizedResp.Header {
		for _, vv := range v {
			ctx.Response.Header.Add(k, vv)
		}
	}
	ctx.Response.Header.Set("X-Cache", "HIT")
	ctx.Response.SetBody(BodyData)
	return nil
}

func (p *Proxy) ServeByOrigin(ctx *fasthttp.RequestCtx) error {
	// MISS
	// TODO: process error in one place
	p.o.ServeProxyHTTP(ctx)
	ctx.Response.Header.Set("X-Cache", "MISS")
	return nil
}
