package conf

type Config struct {
	Discover DiscoverConf `yaml:"discover"`
}

type DiscoverConf struct {
	UseType     string   `yaml:"use_type"`
	DialTimeout int      `yaml:"dial_timeout"`
	Endpoints   []string `yaml:"endpoints"`
}
