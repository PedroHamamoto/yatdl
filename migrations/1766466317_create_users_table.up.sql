CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
