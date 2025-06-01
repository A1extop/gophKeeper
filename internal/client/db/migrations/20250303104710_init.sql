-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lockbox (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    description VARCHAR(255),
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    synced_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    UNIQUE (name, user_id)
    );


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lockbox;
-- +goose StatementEnd
