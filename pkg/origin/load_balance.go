package origin

import (
	"errors"
	"github.com/bocchi-the-cache/hitori/pkg/config"
	"math/rand"
)

func SelectRandomNode(origin *config.OriginSource) (string, error) {
	size := len(origin.Nodes)
	if size == 0 {
		return "", errors.New("no vaild origin source")
	}
	x := rand.Intn(size)
	return origin.Nodes[x], nil
}
