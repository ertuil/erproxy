package conf

// Routes is route rules
type Routes struct {
	Default string            `yaml:"default"`
	Route   map[string]string `yaml:"route"`
}

// TLS a tls
type TLS struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// Auth is  a auth
type Auth struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

// InBound a bound
type InBound struct {
	Type    string            `yaml:"type"`
	Addr    string            `yaml:"addr"`
	UDPPort string            `yaml:"udp"`
	Port    string            `yaml:"port"`
	TLS     TLS               `yaml:"tls"`
	Auth    map[string]string `yaml:"auth"`
}

// OutBound is Bound
type OutBound struct {
	Type   string            `yaml:"type"`
	Addr   string            `yaml:"addr"`
	Port   string            `yaml:"port"`
	Auth   map[string]string `yaml:"auth"`
	UseTLS bool              `yaml:"tls"`
}

// Config is a config
type Config struct {
	Log      string              `yaml:"log"`
	InBound  map[string]InBound  `yaml:"in"`
	OutBound map[string]OutBound `yaml:"out"`
	Routes   Routes              `yaml:"routes"`
}
