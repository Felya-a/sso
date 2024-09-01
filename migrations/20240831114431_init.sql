-- +goose Up
-- +goose StatementBegin
create table if not exists public.user
(
    id    serial primary key,
    email text not null unique,
    password text not null
);
create index if not exists idx_email on public.user (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists public.user
-- +goose StatementEnd