package jsonparse

import (
	"os"
	"toml"
)

var EnableTraceBack = true
var ConfigFilepath = "jsonparse.toml" //flag.String("c", "jsonparse.toml", "config file path")

var sysConfig struct {
	Prefixes struct {
		Common []string `toml:"Common"`
		In     []string `toml:"In"`
		Out    []string `toml:"Out"`
	} `toml:"Prefixes"`
}

func loadConfig() {
	_, err := toml.DecodeFile(ConfigFilepath, &sysConfig)
	if err != nil {
		Traceback(err)
		os.Exit(-1)
	}

}
