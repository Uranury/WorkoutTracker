package workout

import "time"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type UpdateSession struct {
	ID            int64      `json:"session_id"`
	PerformedDate *time.Time `json:"performed_date,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	Name          *string    `json:"name,omitempty"`
}
