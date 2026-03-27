-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_summaries DROP COLUMN timeframe_label;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_summaries ADD COLUMN timeframe_label TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd
