-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS project_summaries (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    project_id UUID NOT NULL REFERENCES projects(id),
    summary_markdown TEXT,
    insights_json JSONB,
    summary_level TEXT NOT NULL,
    start_entry_id UUID NOT NULL REFERENCES entries(id),
    end_entry_id UUID NOT NULL REFERENCES entries(id),
    entries_count INTEGER NOT NULL,
    timeframe_label TEXT NOT NULL,
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX project_summaries_project_id_generated_at_desc ON project_summaries(project_id, generated_at DESC);
CREATE INDEX project_summaries_project_id_end_entry_id ON project_summaries(project_id, end_entry_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_summaries;
-- +goose StatementEnd
