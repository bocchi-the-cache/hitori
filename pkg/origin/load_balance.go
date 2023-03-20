package origin

import (
	"errors"
	"math/rand"

	"github.com/bocchi-the-cache/hitori/pkg/config"
)

func SelectRandomNode(origin *config.OriginSource) (string, error) {
	size := len(origin.Nodes)
	if size == 0 {
		return "", errors.New("no valid origin source")
	}
	x := rand.Intn(size)
	return origin.Nodes[x], nil
}
