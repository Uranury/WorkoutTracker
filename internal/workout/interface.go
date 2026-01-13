package workout

import (
	"context"
	"github.com/Uranury/WorkoutTracker/internal/workout/session"
	"time"
)

type Service interface {
	CreateTemplate(ctx context.Context, userID int64, name, description string) (int64, error)
	AddExerciseToTemplate(ctx context.Context, templateID, exerciseID int64, orderIndex, targetSets, targetReps int) (int64, error)
	StartSession(ctx context.Context, userId int64, name string, templateID *int64) (int64, error)
	AddExerciseToSession(ctx context.Context, sessionID, exerciseID int64, orderIndex int) (int64, error)
	SetSessionFinishTime(ctx context.Context, sessionID int64, finishedAt *time.Time) error
	UpdateSession(ctx context.Context, session UpdateSession) error
	RecordSetToSessionExercise(ctx context.Context, sessionExerciseID int64, setNumber, reps int, weight float64, weightUnit session.WeightUnit) (int64, error)
}
