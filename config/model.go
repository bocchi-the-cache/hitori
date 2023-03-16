package config

import "time"

type Duration time.Duration

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var n int64
	if err := unmarshal(&n); err != nil {
		return err
	}
	*d = Duration(time.Duration(n) * time.Millisecond)
	return nil
}

type Config struct {
	Server  Server  `yaml:"server"`
	Log     Log     `yaml:"log"`
	Mapping Mapping `yaml:"mapping"`
	Cache   Cache   `yaml:"cache"`
}

type Server struct {
	Port       int        `yaml:"port"`
	Connection Connection `yaml:"connection"`
}

type Connection struct {
	Client struct {
		MaxConnectionIdleTime Duration `yaml:"max_connection_idle_time_ms"`
	} `yaml:"client"`
	Origin struct {
		ConnectTimeout     Duration `yaml:"connect_timeout_ms"`
		ReadTimeout        Duration `yaml:"read_timeout_ms"`
		ConnectionPoolSize int      `yaml:"connection_pool_size"`
		MaxReadBufSizeKB   int      `yaml:"max_read_buf_size_kb"`
	} `yaml:"origin"`
}

type Log struct {
	Level          int    `yaml:"level"`
	Path           string `yaml:"path"`
	MaxAge         int    `yaml:"max_age"`
	CutDurationMin int    `yaml:"cut_duration_min"`
	CutSize        int    `yaml:"cut_size_mb"`
}

type Mapping struct {
	Domains       []Domain       `yaml:"domains"`
	OriginSources []OriginSource `yaml:"origin_sources"`
}

type Domain struct {
	DomainName  string `yaml:"domain_name"`
	Origins     string `yaml:"origins"`
	CacheConfig struct {
		Enable bool `yaml:"enable"`
		//CacheRules      []struct {
		//	Enable       bool   `yaml:"enable"`
		//	URLPattern   string `yaml:"url_pattern"`
		//	CodePattern  string `yaml:"code_pattern"`
		//	CacheTTL     int    `yaml:"cache_ttl"`
		//} `yaml:"cache_rules"`
		//CacheKeys []struct {
		//	Enable     bool     `yaml:"enable"`
		//	URLPattern string   `yaml:"url_pattern"`
		//	CacheKey   []string `yaml:"cache_key"`
		//} `yaml:"cache_keys"`
	} `yaml:"cache_config"`
}

type OriginSource struct {
	OriginName string   `yaml:"origin_name"`
	Protocol   string   `yaml:"protocol"`
	Nodes      []string `yaml:"nodes"`
}

type Cache struct {
	Enabled         bool `yaml:"enabled"`
	SliceSizeKb     int  `yaml:"slice_size_kb"`
	ExpirationTimeS int  `yaml:"expiration_time_second"`
	//Ram             CacheRam  `yaml:"ram"`
	Disk CacheDisk `yaml:"disk"`
}

type CacheDisk struct {
	//TolerateFailedCount int               `yaml:"tolerate_failed_count"`
	Devices []CacheDiskDevice `yaml:"devices"`
}

type CacheDiskDevice struct {
	Path   string `yaml:"path"`
	SizeMb int    `yaml:"size_mb"`
}

//type CacheRam struct {
//	SizeMb      int     `yaml:"size_mb"`
//	UsageLimit  float64 `yaml:"usage_limit"`
//	Strategy    string  `yaml:"strategy"`
//	ShmKeyStart int     `yaml:"shm_key_start"`
//}
