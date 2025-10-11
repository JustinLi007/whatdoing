-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS snapshot_anime (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  anime_id UUID UNIQUE NOT NULL,
  name TEXT NOT NULL,
  episode INT NOT NULL,
  version INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE snapshot_anime;
-- +goose StatementEnd
