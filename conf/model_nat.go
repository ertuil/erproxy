package conf

type NatServerConf struct {
	InAddr  string `yaml:"inaddr"`
	OutAddr string `yaml:"outaddr"`
	OutPort string `yaml:"outport"`
	InPort  string `yaml:"inport"`
	InTLS   TLS    `yaml:"intls"`
	OutTLS  TLS    `yaml:"outtls"`
	Auth    string `yaml:"auth"`
}

type NatClientConf struct {
	InAddr  string `yaml:"inaddr"`
	OutAddr string `yaml:"outaddr"`
	OutPort string `yaml:"outport"`
	InPort  string `yaml:"inport"`
	InTLS   bool   `yaml:"intls"`
	OutTLS  bool   `yaml:"outtls"`
	Auth    string `yaml:"auth"`
	beat    int    `yaml:"beat"`
}

type Nat struct {
	Server map[string]NatServerConf `yaml:"server"`
	Client map[string]NatClientConf `yaml:"client"`
}
