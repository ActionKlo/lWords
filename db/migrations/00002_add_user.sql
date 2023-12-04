-- +goose Up
INSERT INTO users (
    username, password, created_at
) VALUES (
    'Klo', 123, CURRENT_TIMESTAMP
)

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
