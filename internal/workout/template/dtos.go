package template

import "github.com/Uranury/WorkoutTracker/internal/exercise"

type Details struct {
	TemplateID  int64              `json:"template_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Exercises   []exercise.Summary `json:"exercises"`
}
