package utils

// TODO: utils is ugly. move to specific pkg.

import (
	"github.com/bocchi-the-cache/hitori/pkg/logger"
	"github.com/valyala/fasthttp"
)

func SetCtxErrorWithLog(ctx *fasthttp.RequestCtx, err error, status int) {
	logger.Error(err.Error())
	ctx.Error(err.Error(), status)
}
