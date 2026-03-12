-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects
ADD COLUMN IF NOT EXISTS last_interacted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE INDEX IF NOT EXISTS projects_last_interacted_at_idx ON projects(last_interacted_at);

-- Create or replace the trigger function for entries table
CREATE OR REPLACE FUNCTION update_project_last_interacted()
RETURNS TRIGGER AS $$
BEGIN
    -- For INSERT and UPDATE
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        UPDATE projects 
        SET last_interacted_at = CURRENT_TIMESTAMP 
        WHERE id = NEW.project_id;
        RETURN NEW;
    -- For DELETE
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE projects 
        SET last_interacted_at = CURRENT_TIMESTAMP 
        WHERE id = OLD.project_id;
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for INSERT, UPDATE, and DELETE
DROP TRIGGER IF EXISTS trigger_update_project_last_interacted_insert ON entries;
CREATE TRIGGER trigger_update_project_last_interacted_insert
    AFTER INSERT ON entries
    FOR EACH ROW
    EXECUTE FUNCTION update_project_last_interacted();

DROP TRIGGER IF EXISTS trigger_update_project_last_interacted_update ON entries;
CREATE TRIGGER trigger_update_project_last_interacted_update
    AFTER UPDATE ON entries
    FOR EACH ROW
    EXECUTE FUNCTION update_project_last_interacted();

DROP TRIGGER IF EXISTS trigger_update_project_last_interacted_delete ON entries;
CREATE TRIGGER trigger_update_project_last_interacted_delete
    AFTER DELETE ON entries
    FOR EACH ROW
    EXECUTE FUNCTION update_project_last_interacted();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_project_last_interacted_insert ON entries;
DROP TRIGGER IF EXISTS trigger_update_project_last_interacted_update ON entries;
DROP TRIGGER IF EXISTS trigger_update_project_last_interacted_delete ON entries;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_project_last_interacted();
DROP INDEX IF EXISTS projects_last_interacted_at_idx;
ALTER TABLE projects DROP COLUMN IF EXISTS last_interacted_at;
-- +goose StatementEnd
