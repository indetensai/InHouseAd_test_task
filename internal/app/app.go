package app

import (
	"log"
	"task/internal/config"
	"task/internal/service"
	"time"
)

func Run() {
	walkService := service.NewPingService()
	ticker := time.Tick(time.Minute)
	data := config.ReadConfig()
	for {
		select {
		case <-ticker:
			result := walkService.Ping(data.Links)
			log.Print(result)
		}
	}

}
