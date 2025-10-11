-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS anime_progress (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  episode INT NOT NULL,
  user_id UUID NOT NULL,
  anime_id UUID NOT NULL,
  UNIQUE(user_id, anime_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE anime_progress;
-- +goose StatementEnd
