-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS new_test_db (
    id int PRIMARY KEY,
    username varchar,
    email varchar,
    phone_number varchar

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS new_test_db;
-- +goose StatementEnd
