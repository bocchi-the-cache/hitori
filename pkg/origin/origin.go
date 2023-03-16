package origin

import (
	"fmt"
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"github.com/bocchi-the-cache/hitori/pkg/logger"
	"github.com/valyala/fasthttp"
)

var DefaultOrigin *Origin

func Init(mapCfg *config.Mapping) {
	DefaultOrigin = NewOrigin(mapCfg)
	logger.Info("origin init successfully", mapCfg)
}

func NewOrigin(mapCfg *config.Mapping) *Origin {
	return &Origin{
		c:  &fasthttp.Client{},
		mp: buildOriginMapping(mapCfg),
	}
}

type Origin struct {
	c  *fasthttp.Client
	mp *Mapping
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

	// TODO: using small slice to fetch origin, streaming to client and cache
	oriRequest := fasthttp.AcquireRequest()
	ctx.Request.CopyTo(oriRequest)

	oriRequest.SetRequestURI(fmt.Sprintf("http://%s%s", node, string(ctx.URI().Path())))
	if ori.OriginHost != "" {
		oriRequest.SetHost(ori.OriginHost)
	} else {
		oriRequest.SetHost(string(ctx.URI().Host()))
	}

	resp := fasthttp.AcquireResponse()
	err = o.c.Do(oriRequest, resp)
	if err != nil {
		SetCtxErrorWithLog(ctx, err, fasthttp.StatusServiceUnavailable)
	}

	//copy response header and body to client
	resp.Header.CopyTo(&ctx.Response.Header)
	ctx.Response.SetBody(resp.Body())
}

func SetCtxErrorWithLog(ctx *fasthttp.RequestCtx, err error, status int) {
	logger.Error(err.Error())
	ctx.Error(err.Error(), status)
}
