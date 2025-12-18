package exercise

import (
	"github.com/lib/pq"
	"time"
)

type Exercise struct {
	ID               int64          `json:"id" db:"id"`
	Name             string         `json:"name" db:"name"`
	PrimaryMuscle    string         `json:"primary_muscle" db:"primary_muscle"`
	SecondaryMuscles pq.StringArray `json:"secondary_muscles" db:"secondary_muscles"`
	IsCompound       bool           `json:"is_compound" db:"is_compound"`
	Description      string         `json:"description" db:"description"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
}
