package proxies

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
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

	err := p.ServeMux(ctx)
	if err != nil {
		logger.Errorf("cache serve error: %v", err)
	}

}

// ServeMux a muxer to decide which source(origin/cache) to serve client
func (p *Proxy) ServeMux(ctx *fasthttp.RequestCtx) error {
	// TODO: complete conditions and purge...
	// to find cache
	if string(ctx.Request.Header.Method()) == "PURGE" {
		goto purgeCache
	}
	if string(ctx.Request.Header.Method()) != "GET" {
		goto directOrigin
	}
	return p.ServeByCache(ctx)

directOrigin:
	return p.ServeByOrigin(ctx)

purgeCache:
	return p.PurgeCache(ctx)
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

func (p *Proxy) PurgeCache(ctx *fasthttp.RequestCtx) error {
	//TODO: 404 for not found, 200 for purged
	bodyKey := ctx.Request.URI().String() + "_body"
	err := p.c.Del(bodyKey)
	if err != nil {
		errw := errors.Wrapf(err, "purge cache error, key: %v", bodyKey)
		ctx.Error(errw.Error(), http.StatusInternalServerError)
		return errw
	}
	// TODO: checksum when serve, in case of **header not deleted** while **body deleted**
	headerKey := ctx.Request.URI().String() + "_header"
	err = p.c.Del(headerKey)
	if err != nil {
		errw := errors.Wrapf(err, "purge cache error, key: %v", headerKey)
		ctx.Error(errw.Error(), http.StatusInternalServerError)
		return errw
	}

	ctx.SuccessString("text/plain", fmt.Sprintf("purge cache success, uri: %v", ctx.Request.URI()))
	return nil
}
