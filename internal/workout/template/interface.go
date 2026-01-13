package template

import "context"

type Repository interface {
	CreateTemplate(ctx context.Context, template Template) (int64, error)
	CreateTemplateExercise(ctx context.Context, template Exercise) (int64, error)
	GetTemplateExercises(ctx context.Context, templateID int64) ([]Exercise, error)
	GetTemplateMaxOrderIndex(ctx context.Context, templateID int64) (int, error)
}
