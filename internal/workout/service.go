package workout

import (
	"context"
	"fmt"
	"github.com/Uranury/WorkoutTracker/pkg/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

type Service interface {
}

type service struct {
	repo Repository
	db   *sqlx.DB
}

func NewService(repo Repository, db *sqlx.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) CreateTemplateWithExercises(ctx context.Context, templ *Template) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	repo := NewRepositoryWithTx(tx)
	templateId, err := repo.CreateTemplate(ctx, *templ)
	if err != nil {
		return 0, fmt.Errorf("create template: %w", err)
	}
	for _, exercise := range templ.Exercises {
		exercise.TemplateID = templateId
		if _, err := repo.CreateTemplateExercise(ctx, exercise); err != nil {
			return 0, fmt.Errorf("create exercise: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return templateId, nil
}

func (s *service) StartSession(ctx context.Context, userId int64, name string, templateID *int64) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	repo := NewRepositoryWithTx(tx)

	session := Session{
		UserID:        userId,
		Name:          name,
		TemplateID:    templateID,
		PerformedDate: time.Now(),
		StartedAt:     utils.TimePtr(time.Now()),
	}

	sessionID, err := repo.CreateSession(ctx, session)
	if err != nil {
		return 0, fmt.Errorf("create session: %w", err)
	}

	if templateID != nil {
		templateExercises, err := repo.GetTemplateExercises(ctx, *templateID)
		if err != nil {
			return 0, fmt.Errorf("get template exercises: %w", err)
		}
		for _, te := range templateExercises {
			se := SessionExercise{
				SessionID:  sessionID,
				ExerciseID: te.ExerciseID,
				OrderIndex: te.OrderIndex,
			}
			if _, err := repo.CreateSessionExercise(ctx, se); err != nil {
				return 0, fmt.Errorf("create session exercise: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return sessionID, nil
}

func (s *service) CreateSessionWithExercises(ctx context.Context, session *Session) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	repo := NewRepositoryWithTx(tx)

	sessionId, err := repo.CreateSession(ctx, *session)
	if err != nil {
		return 0, fmt.Errorf("create session: %w", err)
	}
	for _, exercise := range session.Exercises {
		exercise.SessionID = sessionId
		if _, err := repo.CreateSessionExercise(ctx, exercise); err != nil {
			return 0, fmt.Errorf("create session exercise: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return sessionId, nil
}
