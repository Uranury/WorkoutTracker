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

// TODO: Make template creation accept only essential data and populate the rest itself, don't expect a business model from handlers
// TODO: Implement adding exercises to existing templates (same as adding to sessions)
// TODO: Implement actually recording sets to existing session exercises

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

func (s *service) AddExercisesToSession(ctx context.Context, sessionID, exerciseID int64, orderIndex uint) (int64, error) {
	sessionExercise := SessionExercise{
		SessionID:  sessionID,
		ExerciseID: exerciseID,
		OrderIndex: orderIndex,
	}

	if orderIndex == 0 {
		index, err := s.repo.GetMaxOrderIndex(ctx, sessionID)
		if err != nil {
			return 0, fmt.Errorf("get max order index: %w", err)
		}
		sessionExercise.OrderIndex = index + 1
	}

	sessionExerciseID, err := s.repo.CreateSessionExercise(ctx, sessionExercise)
	if err != nil {
		return 0, fmt.Errorf("create session exercise: %w", err)
	}
	return sessionExerciseID, nil
}
