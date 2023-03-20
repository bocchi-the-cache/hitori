package config

import "github.com/anqur/yasch/pkg/types"

type Config struct {
	Server  Server  `yaml:"server" json:"server"`
	Log     Log     `yaml:"log" json:"log"`
	Mapping Mapping `yaml:"mapping" json:"mapping"`
	Cache   Cache   `yaml:"cache" json:"cache"`
}

type Server struct {
	Port       int        `yaml:"port" json:"port"`
	Connection Connection `yaml:"connection" json:"connection"`
}

type Connection struct {
	Client struct {
		MaxConnIdleDuration types.Duration `yaml:"max_conn_idle_duration" json:"max_conn_idle_duration"`
	} `yaml:"client" json:"client"`
	Origin struct {
		ConnectTimeout types.Duration `yaml:"connect_timeout" json:"connect_timeout"`
		ReadTimeout    types.Duration `yaml:"read_timeout" json:"read_timeout"`
		ConnPoolSize   int            `yaml:"conn_pool_size" json:"conn_pool_size"`
		MaxReadBufSize types.Size     `yaml:"max_read_buf_size" json:"max_read_buf_size"`
	} `yaml:"origin" json:"origin"`
}

type Log struct {
	Level       int            `yaml:"level" json:"level"`
	Path        string         `yaml:"path" json:"path"`
	MaxAge      int            `yaml:"max_age" json:"max_age"`
	CutDuration types.Duration `yaml:"cut_duration" json:"cut_duration"`
	CutSize     types.Size     `yaml:"cut_size" json:"cut_size"`
}

type Mapping struct {
	Domains       []Domain       `yaml:"domains" json:"domains"`
	OriginSources []OriginSource `yaml:"origin_sources" json:"origin_sources"`
}

type Domain struct {
	DomainName  string `yaml:"domain_name" json:"domain_name"`
	Origins     string `yaml:"origins" json:"origins"`
	CacheConfig struct {
		Enabled bool `yaml:"enabled" json:"enabled"`
		//CacheRules      []struct {
		//	Enabled       bool   `yaml:"enabled"`
		//	URLPattern   string `yaml:"url_pattern"`
		//	CodePattern  string `yaml:"code_pattern"`
		//	CacheTTL     int    `yaml:"cache_ttl"`
		//} `yaml:"cache_rules"`
		//CacheKeys []struct {
		//	Enabled     bool     `yaml:"enabled"`
		//	URLPattern string   `yaml:"url_pattern"`
		//	CacheKey   []string `yaml:"cache_key"`
		//} `yaml:"cache_keys"`
	} `yaml:"cache_config" json:"cache_config"`
}

type OriginSource struct {
	OriginName string   `yaml:"origin_name" json:"origin_name"`
	Protocol   string   `yaml:"protocol" json:"protocol"`
	OriginHost string   `yaml:"origin_host,omitempty" json:"origin_host,omitempty"`
	Nodes      []string `yaml:"nodes" json:"nodes"`
}

type Cache struct {
	Enabled   bool           `yaml:"enabled" json:"enabled"`
	SliceSize types.Size     `yaml:"slice_size" json:"slice_size"`
	TTL       types.Duration `yaml:"ttl" json:"ttl"`
	//Ram             CacheRam  `yaml:"ram"`
	Disk CacheDisk `yaml:"disk" json:"disk"`
}

type CacheDisk struct {
	//TolerateFailedCount int               `yaml:"tolerate_failed_count"`
	Devices []CacheDiskDevice `yaml:"devices" json:"devices"`
}

type CacheDiskDevice struct {
	Path string     `yaml:"path" json:"path"`
	Size types.Size `yaml:"size" json:"size"`
}

//type CacheRam struct {
//	SizeMb      int     `yaml:"size_mb"`
//	UsageLimit  float64 `yaml:"usage_limit"`
//	Strategy    string  `yaml:"strategy"`
//	ShmKeyStart int     `yaml:"shm_key_start"`
//}
