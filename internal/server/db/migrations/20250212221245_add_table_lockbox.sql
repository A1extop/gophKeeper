-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lockbox
(
    id   SERIAL PRIMARY KEY,
    name        VARCHAR(1000) NOT NULL,
    url  VARCHAR(255),
    username      VARCHAR(1000),
    password VARCHAR(1000),
    description VARCHAR(1200),
    user_id    INT REFERENCES users (user_id) ON DELETE CASCADE,
    created_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lockbox;
-- +goose StatementEnd
