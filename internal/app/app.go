package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task/internal/config"
	"task/internal/controller"
	"task/internal/repository"
	"task/internal/service"
	"time"
)

func Run() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	pingService := service.NewPingService()
	checkService := service.NewCheckService()
	pingTicker := time.Tick(time.Minute)
	config := config.ReadConfig()
	db := repository.ConnectToDB()
	repo := repository.NewRepository(db)
	result := pingService.Ping(config.Links)
	checkService.UpdateData(result)
	server := controller.NewServer(repo, config, checkService)
	go func() {
		if err := server.Body.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	for {
		select {
		case <-pingTicker:
			result := pingService.Ping(config.Links)
			checkService.UpdateData(result)
		case <-sigint:
			counters := checkService.GetCounters()
			repo.SaveStats(counters["specific"], counters["min"], counters["max"])
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			server.Body.Shutdown(ctx)
		}
	}

}
