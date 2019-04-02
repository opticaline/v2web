package core

import (
	"flag"
)

type Config struct {
	ApiHost    string
	ApiPort    uint
	ServerPort uint
}

func (c *Config) Load() (Config, error) {
	flag.StringVar(&c.ApiHost, "ah", "127.0.0.1", "v2ray api server host")
	flag.UintVar(&c.ApiPort, "ap", 1090, "v2ray api server port")
	flag.UintVar(&c.ServerPort, "p", 8080, "web port")

	flag.Parse()
	return *c, nil
}
