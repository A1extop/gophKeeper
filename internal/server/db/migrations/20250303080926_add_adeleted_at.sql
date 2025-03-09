-- +goose Up
-- +goose StatementBegin
ALTER TABLE lockbox ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE lockbox DROP COLUMN deleted_at;
-- +goose StatementEnd
