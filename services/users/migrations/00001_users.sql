-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  email TEXT UNIQUE NOT NULL,
  password_hash BYTEA NOT NULL,
  refresh_token BYTEA,
  expiry TIMESTAMP,
  username TEXT UNIQUE DEFAULT NULL,
  role TEXT NOT NULL DEFAULT 'regular'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
