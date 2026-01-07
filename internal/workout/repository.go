package workout

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateTemplate(ctx context.Context, template Template) error {
	query := `INSERT INTO workout_templates (user_id, name, description) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, template.UserID, template.Name, template.Description)
	return err
}

func (r *repository) CreateSession(ctx context.Context, session Session) error {
	query := `INSERT INTO workout_sessions (user_id, template_id, performed_date, name, notes) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, session.UserID, session.TemplateID, session.PerformedDate, session.Name, session.Notes)
	return err
}

func (r *repository) CreateSet(ctx context.Context, excSet SessionExerciseSet) error {
	query := `INSERT INTO workout_session_sets (session_exercise_id, set_number, reps, weight, weight_unit) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, excSet.SessionExerciseID, excSet.SetNumber, excSet.Reps, excSet.Weight, excSet.WeightUnit)
	return err
}
