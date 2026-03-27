-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_summaries
ADD COLUMN IF NOT EXISTS summary_text TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_summaries
DROP COLUMN IF EXISTS summary_text;
-- +goose StatementEnd
