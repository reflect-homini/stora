-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS entries (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS entries_project_id_idx ON entries(project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS entries_project_id_idx;
DROP TABLE IF EXISTS entries;
-- +goose StatementEnd
