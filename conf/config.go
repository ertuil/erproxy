package conf

import (
	"log"
	"io/ioutil"
	"github.com/go-yaml/yaml"
)

var (
	CC	Config
)

// GetConfig load config from file
func GetConfig(file string) {
	content,err := ioutil.ReadFile(file)
	if err  != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(content, &CC)
	if err != nil {
		log.Fatalln(err)
	}
}