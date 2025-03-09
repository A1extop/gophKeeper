-- +goose Up
-- +goose StatementBegin
ALTER TABLE lockbox
ADD CONSTRAINT unique_lockbox_name_per_user UNIQUE (name, user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lockbox
DROP CONSTRAINT IF EXISTS unique_lockbox_name_per_user;
-- +goose StatementEnd
