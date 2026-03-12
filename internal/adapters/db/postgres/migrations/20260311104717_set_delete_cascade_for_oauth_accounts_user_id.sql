-- +goose Up
-- +goose StatementBegin
ALTER TABLE oauth_accounts
DROP CONSTRAINT IF EXISTS oauth_accounts_user_id_fkey;

ALTER TABLE oauth_accounts
ADD CONSTRAINT oauth_accounts_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE oauth_accounts
DROP CONSTRAINT IF EXISTS oauth_accounts_user_id_fkey;

ALTER TABLE oauth_accounts
ADD CONSTRAINT oauth_accounts_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users(id);
-- +goose StatementEnd
