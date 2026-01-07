DROP TABLE IF EXISTS exercises, workout_templates, workout_template_exercises, workout_sessions, workout_session_exercises, workout_session_sets;
DROP TYPE IF EXISTS weight_unit;
DROP INDEX IF EXISTS idx_workout_templates_user_id, idx_sessions_user_date, idx_template_exercise_template_id, idx_session_exercise_session_id, idx_exercises_name;