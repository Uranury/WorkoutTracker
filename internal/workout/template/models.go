package template

import "time"

type Template struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	Exercises []Exercise `json:"exercises" db:"-"`
}

type Exercise struct {
	ID         int64 `json:"id" db:"id"`
	TemplateID int64 `json:"template_id" db:"template_id"`
	ExerciseID int64 `json:"exercise_id" db:"exercise_id"`

	OrderIndex int `json:"order_index" db:"order_index"`
	TargetSets int `json:"target_sets" db:"target_sets"`
	TargetReps int `json:"target_reps" db:"target_reps"`
}
