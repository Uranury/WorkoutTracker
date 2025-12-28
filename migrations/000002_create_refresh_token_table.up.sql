CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP NULL,
    user_agent TEXT,
    ip VARCHAR(50)
);

CREATE UNIQUE INDEX idx_refresh_token_hash
ON refresh_tokens (token_hash);

CREATE INDEX idx_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_expires_at ON refresh_tokens (expires_at);

-- For cleanup queries
CREATE INDEX idx_revoked_expired
ON refresh_tokens (revoked_at, expires_at);