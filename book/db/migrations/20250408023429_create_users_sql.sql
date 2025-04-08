-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id integer not null primary key,
    user_id text unique not null,
    username text not null,
    password text not null,
    email text unique not null,
    phone text unique not null,
    role text not null,
    last_login datetime,
    last_login_ip text,
    last_login_user_agent text,
    created_at datetime not null,
    updated_at datetime not null,
    deleted_at datetime
);
-- +goose StatementBegin
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS users;
-- +goose StatementBegin
-- +goose StatementEnd