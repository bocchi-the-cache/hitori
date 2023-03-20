package proxy

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/bocchi-the-cache/hitori/pkg/logger"
	"github.com/bocchi-the-cache/hitori/pkg/origin"
)

const ServerToken = "hitori-cache-server"

var DefaultProxy *HttpProxy

func Init(config *config.Config) {
	DefaultProxy = NewHttpProxy(config.Server.Port, origin.DefaultOrigin)
	logger.Info("proxy init successfully", DefaultProxy.ListenAddr)
}

func Serve() error {
	return DefaultProxy.Serve()
}

type HttpProxy struct {
	ListenAddr string
	s          *fasthttp.Server
	o          *origin.Origin
}

func (p *HttpProxy) Serve() error {
	logger.Info("proxy server start", p.ListenAddr)
	return p.s.ListenAndServe(p.ListenAddr)
}

// NewHttpProxy
// TODO: use buildOption to support complex proxy settings
func NewHttpProxy(port int, ori *origin.Origin) *HttpProxy {
	p := new(HttpProxy)
	p.ListenAddr = fmt.Sprintf(":%d", port)
	p.s = &fasthttp.Server{
		Handler: p.ProxyHandler,
		Name:    ServerToken,
	}
	p.o = ori
	return p
}

func (p *HttpProxy) ProxyHandler(ctx *fasthttp.RequestCtx) {
	logger.Debugf("receive client request uri: %v", ctx.Request.URI())
	ctx.Response.Header.Add("X-Proxy", ServerToken)
	p.o.ServeProxyHTTP(ctx)
}
