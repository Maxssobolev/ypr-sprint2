package models

import "time"

type MovieEventPayload struct {
	MovieID     int      `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"`
	UserID      *int     `json:"user_id,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Description string   `json:"description,omitempty"`
}

type UserEventPayload struct {
	UserID   int     `json:"user_id"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Action   string  `json:"action"`
}

type PaymentEventPayload struct {
	PaymentID int     `json:"payment_id"`
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Method    *string `json:"method_type,omitempty"`
}

type EventDetail struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type SuccessResponse struct {
	Status    string      `json:"status"`
	Partition int         `json:"partition"`
	Offset    int64       `json:"offset"`
	Event     EventDetail `json:"event"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
