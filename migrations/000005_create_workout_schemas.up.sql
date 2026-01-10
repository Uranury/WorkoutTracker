CREATE TABLE IF NOT EXISTS exercises (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    primary_muscle VARCHAR(50) NOT NULL,
    secondary_muscles TEXT[] NOT NULL DEFAULT '{}',
    is_compound BOOLEAN NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_templates (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_template_exercises (
    id BIGSERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL REFERENCES workout_templates(id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
    order_index INT NOT NULL CHECK (order_index > 0),
    target_sets INT NOT NULL CHECK (target_sets > 0),
    target_reps INT NOT NULL CHECK (target_reps > 0),

    UNIQUE(template_id, order_index),
    UNIQUE(template_id, exercise_id)
);

CREATE TABLE IF NOT EXISTS workout_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    template_id BIGINT REFERENCES workout_templates(id) ON DELETE SET NULL,
    performed_date DATE NOT NULL,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    name VARCHAR(100) NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (started_at IS NULL OR finished_at IS NULL OR started_at <= finished_at)
);

CREATE TABLE IF NOT EXISTS workout_session_exercises (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES workout_sessions(id) ON DELETE CASCADE,
    exercise_id BIGINT NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
    order_index INT NOT NULL CHECK (order_index > 0),

    UNIQUE(session_id, order_index),
    UNIQUE(session_id, exercise_id)
);

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_type
            WHERE typname = 'weight_unit'
        ) THEN
            CREATE TYPE weight_unit AS ENUM ('kg', 'lbs');
        END IF;
    END
$$;

CREATE TABLE IF NOT EXISTS workout_session_sets (
    id BIGSERIAL PRIMARY KEY,
    session_exercise_id BIGINT NOT NULL REFERENCES workout_session_exercises(id) ON DELETE CASCADE,
    set_number INT NOT NULL CHECK(set_number > 0),
    reps INT NOT NULL CHECK(reps > 0),
    weight DOUBLE PRECISION NOT NULL CHECK (weight >= 0),
    weight_unit weight_unit NOT NULL,

    UNIQUE(session_exercise_id, set_number)
);

CREATE INDEX idx_workout_templates_user_id ON workout_templates(user_id);
CREATE INDEX idx_sessions_user_date ON workout_sessions(user_id, performed_date DESC);
CREATE INDEX idx_template_exercise_template_id ON workout_template_exercises(template_id);
CREATE INDEX idx_session_exercise_session_id ON workout_session_exercises(session_id);
CREATE INDEX idx_exercises_name ON exercises USING gin(to_tsvector('english', name));