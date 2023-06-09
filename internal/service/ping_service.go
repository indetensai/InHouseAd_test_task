package service

import (
	"net/http"
	"sync"
	"task/internal/models"
	"time"
)

func check_url(url string) models.WebsiteCheck {
	start := time.Now()
	client := http.Client{Timeout: time.Minute}
	r, err := client.Get("http://" + url)
	end := time.Now()
	return models.WebsiteCheck{
		Access:       err == nil && r.StatusCode == 200,
		ResponseTime: end.Sub(start),
	}
}

type PingService interface {
	Ping([]string) models.PingResult
}

type pingService struct {
}

func NewPingService() PingService {
	return &pingService{}
}

func (s *pingService) Ping(links []string) (result models.PingResult) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	result.Data = make(map[string]models.WebsiteCheck)
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			check := check_url(link)
			mu.Lock()
			defer mu.Unlock()
			result.Data[link] = check
		}(link)
	}
	wg.Wait()
	min := result.Data[links[0]].ResponseTime
	max := min
	for url, check := range result.Data {
		if check.ResponseTime >= max {
			result.Max = url
		}
		if check.ResponseTime <= min {
			result.Min = url
		}

	}
	return
}
