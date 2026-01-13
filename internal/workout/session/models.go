package session

import (
	"time"
)

type Session struct {
	ID            int64      `json:"id" db:"id"`
	UserID        int64      `json:"user_id" db:"user_id"`
	TemplateID    *int64     `json:"template_id" db:"template_id"`
	PerformedDate time.Time  `json:"performed_date" db:"performed_date"`
	StartedAt     *time.Time `json:"started_at" db:"started_at"`
	FinishedAt    *time.Time `json:"finished_at" db:"finished_at"`
	Name          string     `json:"name" db:"name"`   // ← "Push Day A", "Legs", etc.
	Notes         string     `json:"notes" db:"notes"` // ← "Felt tired", "New gym"
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	Exercises     []Exercise `json:"exercises" db:"-"` // ← Won't map directly from DB
}

type Exercise struct {
	ID         int64 `json:"id" db:"id"`
	SessionID  int64 `json:"session_id" db:"session_id"`
	ExerciseID int64 `json:"exercise_id" db:"exercise_id"`
	OrderIndex int   `json:"order_index" db:"order_index"`
}

type ExerciseSet struct {
	ID                int64      `json:"id" db:"id"`
	SessionExerciseID int64      `json:"session_exercise_id" db:"session_exercise_id"`
	SetNumber         int        `json:"set_number" db:"set_number"`
	Reps              int        `json:"reps" db:"reps"`
	Weight            float64    `json:"weight" db:"weight"`
	WeightUnit        WeightUnit `json:"weight_unit" db:"weight_unit"`
}

type WeightUnit string

var (
	Kilograms WeightUnit = "kg"
	Pounds    WeightUnit = "lbs"
)
