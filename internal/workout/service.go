package workout

import (
	"context"
	"fmt"
	"github.com/Uranury/WorkoutTracker/pkg/utils"
	"github.com/jmoiron/sqlx"
	"time"
)

type Service interface {
	CreateTemplate(ctx context.Context, userID int64, name, description string) (int64, error)
	AddExerciseToTemplate(ctx context.Context, templateID, exerciseID int64, orderIndex, targetSets, targetReps int) (int64, error)
	StartSession(ctx context.Context, userId int64, name string, templateID *int64) (int64, error)
	AddExerciseToSession(ctx context.Context, sessionID, exerciseID int64, orderIndex int) (int64, error)
	SetSessionFinishTime(ctx context.Context, sessionID int64, finishedAt *time.Time) error
	UpdateSession(ctx context.Context, session UpdateSession) error
	RecordSetToSessionExercise(ctx context.Context, sessionExerciseID int64, setNumber, reps int, weight float64, weightUnit WeightUnit) (int64, error)
}

type service struct {
	repo Repository
	db   *sqlx.DB
}

func NewService(repo Repository, db *sqlx.DB) Service {
	return &service{repo: repo, db: db}
}

func (s *service) CreateTemplate(ctx context.Context, userID int64, name, description string) (int64, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	repo := NewRepositoryWithTx(tx)
	newTemplate := &Template{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	templateId, err := repo.CreateTemplate(ctx, *newTemplate)
	if err != nil {
		return 0, fmt.Errorf("create template: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return templateId, nil
}

func (s *service) AddExerciseToTemplate(ctx context.Context, templateID, exerciseID int64, orderIndex, targetSets, targetReps int) (int64, error) {
	newTemplateExercise := &TemplateExercise{
		TemplateID: templateID,
		ExerciseID: exerciseID,
		OrderIndex: orderIndex,
		TargetSets: targetSets,
		TargetReps: targetReps,
	}

	if orderIndex == 0 {
		currentIndex, err := s.repo.GetTemplateMaxOrderIndex(ctx, templateID)
		if err != nil {
			return 0, fmt.Errorf("get max order index: %w", err)
		}
		newTemplateExercise.OrderIndex = currentIndex + 1
	}

	templateExerciseID, err := s.repo.CreateTemplateExercise(ctx, *newTemplateExercise)
	if err != nil {
		return 0, fmt.Errorf("create exercise: %w", err)
	}
	return templateExerciseID, nil
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

	session := &Session{
		UserID:        userId,
		Name:          name,
		TemplateID:    templateID,
		PerformedDate: time.Now(),
		StartedAt:     utils.TimePtr(time.Now()),
	}

	sessionID, err := repo.CreateSession(ctx, *session)
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

func (s *service) AddExerciseToSession(ctx context.Context, sessionID, exerciseID int64, orderIndex int) (int64, error) {
	sessionExercise := &SessionExercise{
		SessionID:  sessionID,
		ExerciseID: exerciseID,
		OrderIndex: orderIndex,
	}

	if orderIndex == 0 {
		currentIndex, err := s.repo.GetSessionMaxOrderIndex(ctx, sessionID)
		if err != nil {
			return 0, fmt.Errorf("get max order index: %w", err)
		}
		sessionExercise.OrderIndex = currentIndex + 1
	}

	sessionExerciseID, err := s.repo.CreateSessionExercise(ctx, *sessionExercise)
	if err != nil {
		return 0, fmt.Errorf("create session exercise: %w", err)
	}
	return sessionExerciseID, nil
}

func (s *service) SetSessionFinishTime(ctx context.Context, sessionID int64, finishedAt *time.Time) error {
	return s.repo.UpdateSessionFinishTime(ctx, sessionID, finishedAt)
}

func (s *service) UpdateSession(ctx context.Context, session UpdateSession) error {
	return s.repo.UpdateSession(ctx, session.ID, session.Name, session.Notes, session.PerformedDate, session.StartedAt)
}

func (s *service) RecordSetToSessionExercise(ctx context.Context, sessionExerciseID int64, setNumber, reps int, weight float64, weightUnit WeightUnit) (int64, error) {
	performedSet := &SessionExerciseSet{
		SessionExerciseID: sessionExerciseID,
		SetNumber:         setNumber,
		Reps:              reps,
		Weight:            weight,
		WeightUnit:        weightUnit,
	}
	performedSetID, err := s.repo.CreateSet(ctx, *performedSet)
	if err != nil {
		return 0, fmt.Errorf("create exercise set: %w", err)
	}
	return performedSetID, nil
}
