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
	log.Println("startup...")
	config := config.ReadConfig()

	checkService := service.NewCheckService()
	pingService := service.NewPingService()
	db := repository.ConnectToDB()
	repo := repository.NewRepository(db)

	server := controller.NewServer(repo, config, checkService)
	go func() {
		if err := server.Body.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	pingTicker := time.Tick(time.Minute)
	go func() {
		for ; true; <-pingTicker {
			result := pingService.Ping(config.Links)
			checkService.UpdateData(result)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigint
	log.Println("Shutting down server...")

	counters := checkService.GetCounters()
	repo.SaveStats(counters["specific"], counters["min"], counters["max"])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Body.Shutdown(ctx)

	log.Println("Server exiting")
}
