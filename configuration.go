package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	User     string
	Password string
	Server   string
}

func ReadConfig(configfile string) Config {
	var config Config
	source, err := ioutil.ReadFile(configfile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	return config
}
