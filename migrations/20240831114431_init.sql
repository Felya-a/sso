-- +goose Up
-- +goose StatementBegin
create table if not exists chat_app
(
    id    serial primary key,
    email text not null unique
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists chat_app
-- +goose StatementEnd