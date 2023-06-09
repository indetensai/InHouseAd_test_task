package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Links         []string
	ListenAddress string
}

func ReadConfig() Config {
	rawData, err := os.ReadFile("links.txt")
	if err != nil {
		log.Print("failed to read config file")
	}
	data := strings.Split(string(rawData), "\n")
	return Config{
		Links:         data,
		ListenAddress: ":8080",
	}
}
