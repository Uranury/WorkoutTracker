package session

import (
	"context"
	"time"
)

type Repository interface {
	CreateSession(ctx context.Context, session Session) (int64, error)
	CreateSessionExercise(ctx context.Context, session Exercise) (int64, error)
	GetSessionByID(ctx context.Context, sessionID int64) (Session, error)
	GetSessionByTemplateID(ctx context.Context, templateID int64) (Session, error)
	GetSessionMaxOrderIndex(ctx context.Context, sessionID int64) (int, error)
	UpdateSession(ctx context.Context, id int64, name, notes *string, performedDate *time.Time, startedAt *time.Time) error
	UpdateSessionFinishTime(ctx context.Context, sessionID int64, finishedAt *time.Time) error

	CreateSet(ctx context.Context, excSet ExerciseSet) (int64, error)
}
