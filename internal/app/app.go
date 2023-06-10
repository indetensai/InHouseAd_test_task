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
	config := config.ReadConfig()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, os.Kill)
	pingService := service.NewPingService()
	checkService := service.NewCheckService()
	pingTicker := time.Tick(time.Minute)
	db := repository.ConnectToDB()
	repo := repository.NewRepository(db)
	result := pingService.Ping(config.Links)
	checkService.UpdateData(result)
	server := controller.NewServer(repo, config, checkService)
	log.Println("startup...")
	go func() {
		if err := server.Body.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	go func() {
		for {
			select {
			case <-pingTicker:
				result = pingService.Ping(config.Links)
				checkService.UpdateData(result)
			}
		}
	}()
	<-sigint
	log.Println("Shutting down server...")
	counters := checkService.GetCounters()
	repo.SaveStats(counters["specific"], counters["min"], counters["max"])
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Body.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")

}
