package config

import (
	"time"

	"github.com/alecthomas/units"
	"github.com/invopop/jsonschema"
)

type Duration struct{ Dur time.Duration }

func (Duration) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Duration",
		Description: "Go-compatible duration",
	}
}

func (d Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Dur = dur
	return nil
}

type Size struct{ Size units.Base2Bytes }

func (Size) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Size",
		Description: "Go-compatible data size",
	}
}

func (s Size) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var r string
	if err := unmarshal(&r); err != nil {
		return err
	}
	size, err := units.ParseBase2Bytes(r)
	if err != nil {
		return err
	}
	s.Size = size
	return nil
}

type Config struct {
	Server  Server  `yaml:"server" json:"server" jsonschema:"required"`
	Log     Log     `yaml:"log" json:"log" jsonschema:"required"`
	Mapping Mapping `yaml:"mapping" json:"mapping" jsonschema:"required"`
	Cache   Cache   `yaml:"cache" json:"cache" jsonschema:"required"`
}

type Server struct {
	Port       int        `yaml:"port" json:"port" jsonschema:"required"`
	Connection Connection `yaml:"connection" json:"connection" jsonschema:"required"`
}

type Connection struct {
	Client struct {
		MaxConnIdleDuration Duration `yaml:"max_conn_idle_duration" json:"max_conn_idle_duration" jsonschema:"required"`
	} `yaml:"client" json:"client" jsonschema:"required"`
	Origin struct {
		ConnectTimeout Duration `yaml:"connect_timeout" json:"connect_timeout" jsonschema:"required"`
		ReadTimeout    Duration `yaml:"read_timeout" json:"read_timeout" jsonschema:"required"`
		ConnPoolSize   int      `yaml:"conn_pool_size" json:"conn_pool_size" jsonschema:"required"`
		MaxReadBufSize Size     `yaml:"max_read_buf_size" json:"max_read_buf_size" jsonschema:"required"`
	} `yaml:"origin" json:"origin" jsonschema:"required"`
}

type Log struct {
	Level       int      `yaml:"level" json:"level" jsonschema:"required"`
	Path        string   `yaml:"path" json:"path" jsonschema:"required"`
	MaxAge      int      `yaml:"max_age" json:"max_age" jsonschema:"required"`
	CutDuration Duration `yaml:"cut_duration" json:"cut_duration" jsonschema:"required"`
	CutSize     Size     `yaml:"cut_size" json:"cut_size" jsonschema:"required"`
}

type Mapping struct {
	Domains       []Domain       `yaml:"domains" json:"domains" jsonschema:"required"`
	OriginSources []OriginSource `yaml:"origin_sources" json:"origin_sources" jsonschema:"required"`
}

type Domain struct {
	DomainName  string `yaml:"domain_name" json:"domain_name" jsonschema:"required"`
	Origins     string `yaml:"origins" json:"origins" jsonschema:"required"`
	CacheConfig struct {
		Enabled bool `yaml:"enabled" json:"enabled" jsonschema:"required"`
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
	} `yaml:"cache_config" json:"cache_config" jsonschema:"required"`
}

type OriginSource struct {
	OriginName string   `yaml:"origin_name" json:"origin_name" jsonschema:"required"`
	Protocol   string   `yaml:"protocol" json:"protocol" jsonschema:"required"`
	Nodes      []string `yaml:"nodes" json:"nodes" jsonschema:"required"`
}

type Cache struct {
	Enabled   bool     `yaml:"enabled" json:"enabled" jsonschema:"required"`
	SliceSize Size     `yaml:"slice_size" json:"slice_size" jsonschema:"required"`
	TTL       Duration `yaml:"ttl" json:"ttl" jsonschema:"required"`
	//Ram             CacheRam  `yaml:"ram"`
	Disk CacheDisk `yaml:"disk" json:"disk" jsonschema:"required"`
}

type CacheDisk struct {
	//TolerateFailedCount int               `yaml:"tolerate_failed_count"`
	Devices []CacheDiskDevice `yaml:"devices" json:"devices" jsonschema:"required"`
}

type CacheDiskDevice struct {
	Path string `yaml:"path" json:"path" jsonschema:"required"`
	Size Size   `yaml:"size" json:"size" jsonschema:"required"`
}

//type CacheRam struct {
//	SizeMb      int     `yaml:"size_mb"`
//	UsageLimit  float64 `yaml:"usage_limit"`
//	Strategy    string  `yaml:"strategy"`
//	ShmKeyStart int     `yaml:"shm_key_start"`
//}
