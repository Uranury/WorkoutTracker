package workout

import (
	"context"
	"fmt"
	"github.com/Uranury/WorkoutTracker/pkg/database"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository interface {
	CreateTemplate(ctx context.Context, template Template) (int64, error)
	CreateTemplateExercise(ctx context.Context, template TemplateExercise) (int64, error)
	GetTemplateExercises(ctx context.Context, templateID int64) ([]TemplateExercise, error)
	GetTemplateMaxOrderIndex(ctx context.Context, templateID int64) (uint, error)

	CreateSession(ctx context.Context, session Session) (int64, error)
	CreateSessionExercise(ctx context.Context, session SessionExercise) (int64, error)
	GetSessionByID(ctx context.Context, sessionID int64) (Session, error)
	GetSessionByTemplateID(ctx context.Context, templateID int64) (Session, error)
	GetSessionMaxOrderIndex(ctx context.Context, sessionID int64) (uint, error)
	UpdateSession(ctx context.Context, id int64, name, notes *string, performedDate *time.Time, startedAt *time.Time) error
	UpdateSessionFinishTime(ctx context.Context, sessionID int64, finishedAt *time.Time) error

	CreateSet(ctx context.Context, excSet SessionExerciseSet) (int64, error)
}

type repository struct {
	executor database.Executor
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		executor: db,
	}
}

func NewRepositoryWithTx(tx *sqlx.Tx) Repository {
	return &repository{
		executor: tx,
	}
}

func (r *repository) CreateTemplate(ctx context.Context, template Template) (int64, error) {
	query := `INSERT INTO workout_templates (user_id, name, description) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.executor.QueryRowxContext(ctx, query, template.UserID, template.Name, template.Description).Scan(&id)
	return id, err
}

func (r *repository) CreateTemplateExercise(ctx context.Context, te TemplateExercise) (int64, error) {
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

func (r *repository) GetTemplateExercises(ctx context.Context, templateID int64) ([]TemplateExercise, error) {
	var exercises []TemplateExercise
	query := `SELECT * FROM workout_template_exercises WHERE template_id = $1 ORDER BY order_index`
	err := r.executor.SelectContext(ctx, &exercises, query, templateID)
	return exercises, err
}

func (r *repository) GetTemplateMaxOrderIndex(ctx context.Context, templateID int64) (uint, error) {
	query := `SELECT COALESCE(MAX(order_index), 0) FROM workout_template_exercises WHERE template_id = $1`
	var order uint
	if err := r.executor.QueryRowxContext(ctx, query, templateID).Scan(&order); err != nil {
		return 0, err
	}
	return order, nil
}

func (r *repository) CreateSession(ctx context.Context, session Session) (int64, error) {
	query := `INSERT INTO workout_sessions (user_id, template_id, performed_date, name, notes, started_at, finished_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var id int64
	err := r.executor.QueryRowxContext(ctx, query, session.UserID, session.TemplateID, session.PerformedDate, session.Name, session.Notes, session.StartedAt, session.FinishedAt).Scan(&id)
	return id, err
}

func (r *repository) CreateSessionExercise(ctx context.Context, se SessionExercise) (int64, error) {
	query := `INSERT INTO workout_session_exercises (session_id, exercise_id, order_index) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := r.executor.QueryRowxContext(ctx, query, se.SessionID, se.ExerciseID, se.OrderIndex).Scan(&id)
	return id, err
}

func (r *repository) GetSessionByID(ctx context.Context, sessionID int64) (Session, error) {
	var session Session
	query := `SELECT * FROM workout_sessions WHERE id = $1`
	if err := r.executor.GetContext(ctx, &session, query, sessionID); err != nil {
		return Session{}, err
	}
	return session, nil
}

func (r *repository) GetSessionMaxOrderIndex(ctx context.Context, sessionID int64) (uint, error) {
	query := `SELECT COALESCE(MAX(order_index), 0) FROM workout_session_exercises WHERE session_id = $1`
	var order uint
	if err := r.executor.QueryRowxContext(ctx, query, sessionID).Scan(&order); err != nil {
		return 0, err
	}
	return order, nil
}

func (r *repository) GetSessionByTemplateID(ctx context.Context, templateID int64) (Session, error) {
	query := `SELECT * FROM workout_sessions WHERE template_id = $1`
	var session Session
	if err := r.executor.QueryRowxContext(ctx, query, templateID).StructScan(&session); err != nil {
		return Session{}, err
	}
	return session, nil
}

func (r *repository) UpdateSession(ctx context.Context, id int64, name, notes *string, performedDate *time.Time, startedAt *time.Time) error {
	query := `UPDATE workout_sessions
              SET 
              name = COALESCE($1, name),
              notes = COALESCE($2, notes),
              performed_date = COALESCE($3, performed_date),
              started_at = COALESCE($4, started_at)
              WHERE id = $5`
	_, err := r.executor.ExecContext(ctx, query, name, notes, performedDate, startedAt, id)
	return err
}

func (r *repository) UpdateSessionFinishTime(ctx context.Context, id int64, finishedAt *time.Time) error {
	query := `UPDATE workout_sessions SET finished_at = $1 WHERE id = $2`

	res, err := r.executor.ExecContext(ctx, query, finishedAt, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("session %d not found", id)
	}
	return nil
}

func (r *repository) CreateSet(ctx context.Context, excSet SessionExerciseSet) (int64, error) {
	query := `INSERT INTO workout_session_sets (session_exercise_id, set_number, reps, weight, weight_unit) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := r.executor.QueryRowxContext(ctx, query, excSet.SessionExerciseID, excSet.SetNumber, excSet.Reps, excSet.Weight, excSet.WeightUnit).Scan(&id)
	return id, err
}
