package origin

import "github.com/bocchi-the-cache/hitori/pkg/config"

type Mapping struct {
	config.Mapping
	DomainMap map[string]*config.Domain
	OriginMap map[string]*config.OriginSource
}

func buildOriginMapping(mapCfg *config.Mapping) *Mapping {
	mp := &Mapping{
		Mapping:   *mapCfg,
		DomainMap: make(map[string]*config.Domain),
		OriginMap: make(map[string]*config.OriginSource),
	}
	for _, domain := range mapCfg.Domains {
		mp.DomainMap[domain.DomainName] = &domain
	}
	for _, source := range mapCfg.OriginSources {
		mp.OriginMap[source.OriginName] = &source
	}
	return mp
}
