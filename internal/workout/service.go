package workout

import (
	"context"
	"fmt"
	"github.com/Uranury/WorkoutTracker/internal/workout/session"
	"github.com/Uranury/WorkoutTracker/internal/workout/template"
	"github.com/Uranury/WorkoutTracker/pkg/database"
	"github.com/Uranury/WorkoutTracker/pkg/utils"
	"time"
)

type service struct {
	templateRepo template.Repository
	sessionRepo  session.Repository
	txProvider   database.TxProvider
}

func NewService(templateRepo template.Repository, sessionRepo session.Repository, txProvider database.TxProvider) Service {
	return &service{templateRepo, sessionRepo, txProvider}
}

func (s *service) CreateTemplate(ctx context.Context, userID int64, name, description string) (int64, error) {
	newTemplate := &template.Template{
		UserID:      userID,
		Name:        name,
		Description: description,
	}
	templateId, err := s.templateRepo.CreateTemplate(ctx, *newTemplate)
	if err != nil {
		return 0, fmt.Errorf("create template: %w", err)
	}
	return templateId, nil
}

func (s *service) UpdateTemplate(ctx context.Context, templateID int64, name, description string) (template.Template, error) {
	return s.templateRepo.UpdateTemplate(ctx, templateID, name, description)
}

func (s *service) DeleteTemplate(ctx context.Context, templateID int64) error {
	return s.templateRepo.DeleteTemplate(ctx, templateID)
}

func (s *service) AddExerciseToTemplate(ctx context.Context, templateID, exerciseID int64, orderIndex, targetSets, targetReps int) (int64, error) {
	newTemplateExercise := &template.Exercise{
		TemplateID: templateID,
		ExerciseID: exerciseID,
		OrderIndex: orderIndex,
		TargetSets: targetSets,
		TargetReps: targetReps,
	}

	if orderIndex == 0 {
		currentIndex, err := s.templateRepo.GetTemplateMaxOrderIndex(ctx, templateID)
		if err != nil {
			return 0, fmt.Errorf("get max order index: %w", err)
		}
		newTemplateExercise.OrderIndex = currentIndex + 1
	}

	templateExerciseID, err := s.templateRepo.CreateTemplateExercise(ctx, *newTemplateExercise)
	if err != nil {
		return 0, fmt.Errorf("create exercise: %w", err)
	}
	return templateExerciseID, nil
}

func (s *service) StartSession(ctx context.Context, userId int64, name string, templateID *int64) (int64, error) {
	var sessionID int64

	err := s.txProvider.RunInTx(ctx, func(exec database.Executor) error {
		sessRepo := session.NewRepository(exec)
		tmplRepo := template.NewRepository(exec)

		newSession := &session.Session{
			UserID:        userId,
			Name:          name,
			TemplateID:    templateID,
			PerformedDate: time.Now(),
			StartedAt:     utils.TimePtr(time.Now()),
		}

		sessionID, err := sessRepo.CreateSession(ctx, *newSession)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		if templateID != nil {
			templateExercises, err := tmplRepo.GetTemplateExercises(ctx, *templateID)
			if err != nil {
				return fmt.Errorf("get template exercises: %w", err)
			}
			for _, te := range templateExercises {
				se := session.Exercise{
					SessionID:  sessionID,
					ExerciseID: te.ExerciseID,
					OrderIndex: te.OrderIndex,
				}
				if _, err := sessRepo.CreateSessionExercise(ctx, se); err != nil {
					return fmt.Errorf("create session exercise: %w", err)
				}
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return sessionID, nil
}

func (s *service) AddExerciseToSession(ctx context.Context, sessionID, exerciseID int64, orderIndex int) (int64, error) {
	sessionExercise := &session.Exercise{
		SessionID:  sessionID,
		ExerciseID: exerciseID,
		OrderIndex: orderIndex,
	}

	if orderIndex == 0 {
		currentIndex, err := s.sessionRepo.GetSessionMaxOrderIndex(ctx, sessionID)
		if err != nil {
			return 0, fmt.Errorf("get max order index: %w", err)
		}
		sessionExercise.OrderIndex = currentIndex + 1
	}

	sessionExerciseID, err := s.sessionRepo.CreateSessionExercise(ctx, *sessionExercise)
	if err != nil {
		return 0, fmt.Errorf("create session exercise: %w", err)
	}
	return sessionExerciseID, nil
}

func (s *service) SetSessionFinishTime(ctx context.Context, sessionID int64, finishedAt *time.Time) error {
	return s.sessionRepo.UpdateSessionFinishTime(ctx, sessionID, finishedAt)
}

func (s *service) UpdateSession(ctx context.Context, session UpdateSession) error {
	return s.sessionRepo.UpdateSession(ctx, session.ID, session.Name, session.Notes, session.PerformedDate, session.StartedAt)
}

func (s *service) RecordSetToSessionExercise(ctx context.Context, sessionExerciseID int64, setNumber, reps int, weight float64, weightUnit session.WeightUnit) (int64, error) {
	performedSet := &session.ExerciseSet{
		SessionExerciseID: sessionExerciseID,
		SetNumber:         setNumber,
		Reps:              reps,
		Weight:            weight,
		WeightUnit:        weightUnit,
	}
	performedSetID, err := s.sessionRepo.CreateSet(ctx, *performedSet)
	if err != nil {
		return 0, fmt.Errorf("create exercise set: %w", err)
	}
	return performedSetID, nil
}
