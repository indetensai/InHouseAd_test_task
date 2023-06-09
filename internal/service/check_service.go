package service

import (
	"fmt"
	"sync"
	"sync/atomic"
	"task/internal/models"
	"task/internal/repository"
	"time"
)

type CheckService interface {
	SpecificCheck(link string) (time.Duration, error)
	MinResponseTime() string
	MaxResponseTime() string
	UpdateData(models.PingResult)
	CountersUpdate(data map[string]uint64)
	SpecificCounterIncrementing()
	SlowestCounterIncrementing()
	FastestCounterIncrementing()
	GetCounters() map[string]uint64
}

type checkService struct {
	repo            repository.Repository
	data            models.PingResult
	mu              sync.RWMutex
	SpecificCounter uint64
	SlowestCounter  uint64
	FastestCounter  uint64
}

func NewCheckService() CheckService {
	return &checkService{data: models.PingResult{}}
}

func (c *checkService) UpdateData(input models.PingResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = input
}

func (c *checkService) SpecificCheck(link string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if _, ok := c.data.Data[link]; !ok {
		return 0, fmt.Errorf("link is not in map")
	}
	return c.data.Data[link].ResponseTime, nil
}

func (c *checkService) MinResponseTime() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data.Min
}

func (c *checkService) MaxResponseTime() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data.Max
}

func (c *checkService) SpecificCounterIncrementing() {
	atomic.AddUint64(&c.SpecificCounter, 1)
}

func (c *checkService) SlowestCounterIncrementing() {
	atomic.AddUint64(&c.SlowestCounter, 1)
}

func (c *checkService) FastestCounterIncrementing() {
	atomic.AddUint64(&c.FastestCounter, 1)
}

func (c *checkService) CountersUpdate(data map[string]uint64) {
	c.SpecificCounter = data["specific"]
	c.SlowestCounter = data["max"]
	c.FastestCounter = data["min"]
}

func (c *checkService) GetCounters() map[string]uint64 {
	result := make(map[string]uint64)
	result["specific"] = c.SpecificCounter
	result["max"] = c.SlowestCounter
	result["min"] = c.FastestCounter
	return result
}
