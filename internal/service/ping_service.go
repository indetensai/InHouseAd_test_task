package service

import (
	"net/http"
	"sync"
	"task/internal/models"
	"task/internal/repository"
	"time"
)

func url_ok(url string) models.WebsiteCheck {
	start := time.Now()
	r, err := http.Head(url)
	end := time.Now()
	return models.WebsiteCheck{
		Access:       err == nil && r.StatusCode == 200,
		ResponseTime: end.Sub(start),
	}
}

type PingService interface {
	Ping([]string) sync.Map
}

type pingService struct {
	repo repository.Repository
}

func NewPingService() PingService {
	return &pingService{}
}

func (s *pingService) Ping(links []string) sync.Map {
	var wg sync.WaitGroup
	var result sync.Map
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			result.Store(link, url_ok(link))
		}(link)
	}
	wg.Wait()
	return result
}
