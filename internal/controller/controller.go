package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"task/internal/repository"
	"task/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
)

type responseTimeHandler struct {
	repo          repository.Repository
	check_service service.CheckService
}

func AddHandlers(r chi.Router, service service.CheckService, repo repository.Repository) {
	saveTicker := time.Tick(time.Hour)

	handler := responseTimeHandler{check_service: service, repo: repo}
	data, err := handler.repo.GetStats()
	if err != nil {
		log.Print("failed to get stats")
	}
	if data != nil {
		handler.check_service.CountersUpdate(data)
	}
	r.Get("/response-time/of", handler.SpecificResponseTime)
	r.Get("/response-time/slowest", handler.SlowestResponseTime)
	r.Get("/response-time/fastest", handler.FastestResponseTime)
	r.Get("/response-time/stats", handler.Stats)
	go func() {
		for {
			select {
			case <-saveTicker:
				counters := handler.check_service.GetCounters()
				handler.repo.SaveStats(counters["specific"], counters["min"], counters["max"])
			}
		}
	}()
}

func (h *responseTimeHandler) SpecificResponseTime(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.check_service.SpecificCounterIncrementing()
	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	responseTime, err := h.check_service.SpecificCheck(url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	result, err := json.Marshal(map[string]time.Duration{"response time": time.Duration(responseTime.Milliseconds())})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func (h *responseTimeHandler) SlowestResponseTime(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.check_service.SlowestCounterIncrementing()
	url := h.check_service.MaxResponseTime()
	result, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func (h *responseTimeHandler) FastestResponseTime(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.check_service.FastestCounterIncrementing()
	url := h.check_service.MinResponseTime()
	result, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func (h *responseTimeHandler) Stats(
	w http.ResponseWriter,
	r *http.Request) {
	counters := h.check_service.GetCounters()
	result, err := json.Marshal(map[string]uint64{
		"/response-time/of":      counters["specific"],
		"/response-time/slowest": counters["max"],
		"/response-time/fastest": counters["min"],
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}
