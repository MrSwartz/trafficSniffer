package config

import (
	"gopkg.in/gcfg.v1"
	e "saas/internal/debug/err"
)

type Config struct {
	Ram struct {
		MaxRam   uint64
		Priority uint8
		Enabled  bool
	}
	Scanner struct {
		StartListenPort uint32
		EndListenPort   uint32
		Priority        uint8
		Enabled         bool
	}
	Location struct {
		IPv4      string
		IPv6      string
		ProxyIPv4 string
		ProxyIPv6 string
		Priority  uint8
		Enabled   bool
	}
	Manufacture struct {
		Mac      string
		Priority uint8
		Enabled  bool
	}
}

func ReadConfig(config string) Config {
	var conf Config
	err := gcfg.ReadFileInto(&conf, config)
	e.CheckErr(err)
	return conf
}
