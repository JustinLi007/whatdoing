-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  anime_id UUID NOT NULL,
  name TEXT NOT NULL,
  episodes INT NOT NULL,
  event_type TEXT NOT NULL,
  status TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbox;
-- +goose StatementEnd
