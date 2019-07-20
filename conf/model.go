package conf

// TLS a tls
type TLS struct {
	Cert 		string		`yaml:"cert"`
	Key 		string		`yaml:"key"`
}

// Auth is  a auth
type Auth struct {
	User		string		`yaml:"user"`
	Token		string		`yaml:"token"`
}

// InBound a bound
type InBound struct {
	Addr		string		`yaml:"addr"`
	Port		string		`yaml:"port"`
	TLS			TLS			`yaml:"tls"`
	Auth		Auth		`yaml:"auth"`
}

// OutBound is Bound
type OutBound struct {
	Type		string		`yaml:"type"`
	Addr		string		`yaml:"addr"`
	Port		string		`yaml:"port"`
	Auth		Auth		`yaml:"auth"`
	UseTLS		bool		`yaml:"tls"`
}
// Config is a config
type Config struct {
	InBound 	InBound		`yaml:"in"`
	OutBound 	OutBound	`yaml:"out"`
}