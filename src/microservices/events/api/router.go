package api

import (
	"net/http"

	"events/kafka"
)

func SetupRouter(producer *kafka.Producer) *http.ServeMux {
	handler := &EventHandler{
		producer: producer,
	}

	router := http.NewServeMux()

	router.HandleFunc("POST /api/events/movie", handler.CreateMovieEvent)
	router.HandleFunc("POST /api/events/user", handler.CreateUserEvent)
	router.HandleFunc("POST /api/events/payment", handler.CreatePaymentEvent)

	router.HandleFunc("GET /api/events/health", handler.HealthCheck)

	return router
}
