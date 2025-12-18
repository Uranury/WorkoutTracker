package workout

import "time"

type Session struct {
	ID        int64      `json:"id" db:"id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	Date      time.Time  `json:"date" db:"date"`
	Name      string     `json:"name" db:"name"`   // ← "Push Day A", "Legs", etc.
	Notes     string     `json:"notes" db:"notes"` // ← "Felt tired", "New gym"
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	Exercises []Exercise `json:"exercises" db:"-"` // ← Won't map directly from DB
}

type Exercise struct {
	ID         int64     `json:"id" db:"id"`
	WorkoutID  int64     `json:"workout_id" db:"workout_id"`
	ExerciseID int64     `json:"exercise_id" db:"exercise_id"`
	Sets       int       `json:"sets" db:"sets"`
	Reps       int       `json:"reps" db:"reps"`
	Weight     float64   `json:"weight" db:"weight"`
	WeightUnit string    `json:"weight_unit" db:"weight_unit"` // ← "kg" or "lbs"
	Notes      string    `json:"notes" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
