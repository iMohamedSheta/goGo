-- +goose Up
CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT NULL,
    status TINYINT DEFAULT 0 NOT NULL, -- in app/enums/todo_status_enum.go
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    INDEX (user_id)
);

-- +goose Down
DROP TABLE IF EXISTS todos;