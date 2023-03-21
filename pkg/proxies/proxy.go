package proxies

import (
	"github.com/bocchi-the-cache/hitori/internal/proxies"
)

var NewProxy = proxies.NewProxy

type Proxy interface{ Serve() error }
