package origin

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/bocchi-the-cache/hitori/pkg/cache"
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/bocchi-the-cache/hitori/pkg/logger"
)

var DefaultOrigin *Origin

func Init(mapCfg *config.Mapping) {
	DefaultOrigin = NewOrigin(mapCfg, cache.DefaultCache)
	logger.Info("origin init successfully", mapCfg)
}

func NewOrigin(mapCfg *config.Mapping, ca cache.Cache) *Origin {
	return &Origin{
		c:     new(fasthttp.Client),
		cache: ca,
		mp:    buildOriginMapping(mapCfg),
	}
}

type Origin struct {
	c     *fasthttp.Client
	cache cache.Cache
	mp    *Mapping
}

// ServeProxyHTTP
// TODO: high performance tuning with `fasthttp`
func (o *Origin) ServeProxyHTTP(ctx *fasthttp.RequestCtx) {
	u := ctx.URI()
	d, ok := o.mp.DomainMap[string(u.Host())]
	if !ok {
		SetCtxErrorWithLog(ctx, fmt.Errorf("proxy domain not found: %v", string(u.Host())), fasthttp.StatusInternalServerError)
		return
	}
	ori, ok := o.mp.OriginMap[d.Origins]
	if !ok {
		SetCtxErrorWithLog(ctx, fmt.Errorf("origin name not found: %v", string(d.Origins)), fasthttp.StatusInternalServerError)
		return
	}
	node, err := SelectRandomNode(ori)
	if err != nil {
		SetCtxErrorWithLog(ctx, err, fasthttp.StatusInternalServerError)
		return
	}
	logger.Debugf("origin select, origin node: %v, uri: %v", ori.OriginName+node, u)

	// TODO: using small slice to fetch origin, streaming to client and cache
	oriRequest := fasthttp.AcquireRequest()
	ctx.Request.CopyTo(oriRequest)

	oriRequest.SetRequestURI(fmt.Sprintf("%s://%s%s", ori.Protocol, node, string(ctx.URI().Path())))
	if ori.OriginHost != "" {
		oriRequest.SetHost(ori.OriginHost)
	} else {
		oriRequest.SetHost(string(ctx.URI().Host()))
	}

	resp := fasthttp.AcquireResponse()
	err = o.c.Do(oriRequest, resp)
	if err != nil {
		SetCtxErrorWithLog(ctx, err, fasthttp.StatusServiceUnavailable)
		return
	}
	logger.Debugf("do origin finished, header: %v, uri: %v", string(resp.Header.Header()), u)

	resp.Header.Set("X-Cache-Timestamp", fmt.Sprintf("%v", time.Now()))
	header := resp.Header.Header()
	body := resp.Body()

	// copy response header and body to cache
	ProduceCache(ctx, o.cache, header, body)

	// copy response header and body to client
	resp.Header.CopyTo(&ctx.Response.Header)
	ctx.Response.SetBody(body)
}

func SetCtxErrorWithLog(ctx *fasthttp.RequestCtx, err error, status int) {
	logger.Error(err.Error())
	ctx.Error(err.Error(), status)
}

// ProduceCache TODO: process header/body in pkg `cache`
func ProduceCache(ctx *fasthttp.RequestCtx, ca cache.Cache, header []byte, body []byte) {
	cacheKey := ctx.Request.URI().String() + "_body"
	err := ca.Set(cacheKey, body)
	if err != nil {
		logger.Error("cache body set error, uri", err)
		return
	}
	logger.Debugf("cache body set, uri: %v", ctx.Request.URI().String())

	// set header cache after body cache finished
	cacheKey = ctx.Request.URI().String() + "_header"
	err = ca.Set(cacheKey, header)
	if err != nil {
		logger.Error("cache header set error, uri", err)
	}
	logger.Debugf("cache header set, uri: %v", ctx.Request.URI().String())
}
