package template

import (
	"context"
	"github.com/Uranury/WorkoutTracker/pkg/database"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	executor database.Executor
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{executor: db}
}

func NewRepositoryWithTx(tx *sqlx.Tx) Repository {
	return &repository{executor: tx}
}

func (r *repository) CreateTemplate(ctx context.Context, template Template) (int64, error) {
	query := `INSERT INTO workout_templates (user_id, name, description) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.executor.QueryRowxContext(ctx, query, template.UserID, template.Name, template.Description).Scan(&id)
	return id, err
}

func (r *repository) CreateTemplateExercise(ctx context.Context, te Exercise) (int64, error) {
	query := `INSERT INTO workout_template_exercises 
              (template_id, exercise_id, order_index, target_sets, target_reps) 
              VALUES ($1, $2, $3, $4, $5) 
              RETURNING id`

	var id int64
	err := r.executor.QueryRowxContext(ctx, query,
		te.TemplateID,
		te.ExerciseID,
		te.OrderIndex,
		te.TargetSets,
		te.TargetReps,
	).Scan(&id)

	return id, err
}

func (r *repository) GetTemplateExercises(ctx context.Context, templateID int64) ([]Exercise, error) {
	var exercises []Exercise
	query := `SELECT * FROM workout_template_exercises WHERE template_id = $1 ORDER BY order_index`
	err := r.executor.SelectContext(ctx, &exercises, query, templateID)
	return exercises, err
}

func (r *repository) GetTemplateMaxOrderIndex(ctx context.Context, templateID int64) (int, error) {
	query := `SELECT COALESCE(MAX(order_index), 0) FROM workout_template_exercises WHERE template_id = $1`
	var order int
	if err := r.executor.QueryRowxContext(ctx, query, templateID).Scan(&order); err != nil {
		return 0, err
	}
	return order, nil
}
