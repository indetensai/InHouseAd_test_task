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
		log.Fatal("failed to read config file")
	}
	data := strings.ReplaceAll(string(rawData), "\r\n", "\n")
	links := strings.Split(data, "\n")
	return Config{
		Links:         links,
		ListenAddress: ":8080",
	}
}
