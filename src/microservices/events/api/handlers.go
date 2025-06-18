package api

import (
	"encoding/json"
	"events/kafka"
	"events/models"
	"net/http"
	"strconv"
	"time"
)

type EventHandler struct {
	producer *kafka.Producer
}

func (h *EventHandler) CreateMovieEvent(w http.ResponseWriter, r *http.Request) {
	var payload models.MovieEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	event := models.EventDetail{
		ID:        "movie-" + strconv.Itoa(payload.MovieID) + "-" + payload.Action,
		Type:      "movie",
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	msgValue, _ := json.Marshal(event)
	result, err := h.producer.SendMessage(r.Context(), "movie-events", []byte("movie"), msgValue)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send event to Kafka")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.SuccessResponse{
		Status:    "success",
		Partition: result.Partition,
		Offset:    result.Offset,
		Event:     event,
	})
}

func (h *EventHandler) CreateUserEvent(w http.ResponseWriter, r *http.Request) {
	var payload models.UserEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	event := models.EventDetail{
		ID:        "user-" + strconv.Itoa(payload.UserID) + "-" + payload.Action,
		Type:      "user",
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	msgValue, _ := json.Marshal(event)
	result, err := h.producer.SendMessage(r.Context(), "user-events", []byte("user"), msgValue)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send event to Kafka")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.SuccessResponse{
		Status:    "success",
		Partition: result.Partition,
		Offset:    result.Offset,
		Event:     event,
	})
}

func (h *EventHandler) CreatePaymentEvent(w http.ResponseWriter, r *http.Request) {
	var payload models.PaymentEventPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	event := models.EventDetail{
		ID:        "payment-" + strconv.Itoa(payload.UserID) + "-" + payload.Status,
		Type:      "payment",
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	msgValue, _ := json.Marshal(event)
	result, err := h.producer.SendMessage(r.Context(), "payment-events", []byte("payment"), msgValue)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to send event to Kafka")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.SuccessResponse{
		Status:    "success",
		Partition: result.Partition,
		Offset:    result.Offset,
		Event:     event,
	})
}

func (h *EventHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]bool{"status": true})
}
