-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS project_summaries_project_id_summary_level_start_entry_id_end_entry_id_uniq_idx
ON project_summaries (project_id, summary_level, start_entry_id, end_entry_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS project_summaries_project_id_summary_level_start_entry_id_end_entry_id_uniq_idx;
-- +goose StatementEnd
