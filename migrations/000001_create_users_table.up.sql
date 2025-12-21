CREATE TABLE IF NOT EXISTS users (
         id BIGSERIAL PRIMARY KEY,
         username VARCHAR(50) UNIQUE NOT NULL,
         email VARCHAR(255) UNIQUE NOT NULL,
         age INT NOT NULL CHECK (age >= 13 AND age <= 120),
         gender VARCHAR(20) NOT NULL CHECK (gender IN ('male', 'female')),
         password VARCHAR(255) NOT NULL,
         created_at TIMESTAMP DEFAULT NOW(),
         updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);