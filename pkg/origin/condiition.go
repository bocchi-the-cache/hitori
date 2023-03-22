package origin

import "github.com/valyala/fasthttp"

func isResponseOKToCache(resp *fasthttp.Response) bool {
	if resp.StatusCode() == fasthttp.StatusOK {
		return true
	}
	return false
}
