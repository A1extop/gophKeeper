-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_type AS ENUM ('admin', 'attendee');

CREATE TABLE IF NOT EXISTS Users
(
    user_id       SERIAL PRIMARY KEY,
    username      VARCHAR(50)         NOT NULL UNIQUE,
    password_hash VARCHAR(255)        NOT NULL,
    user_type     user_type           NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users;
DROP TYPE IF EXISTS user_type;
-- +goose StatementEnd
