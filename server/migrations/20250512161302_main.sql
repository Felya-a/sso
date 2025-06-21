-- +goose Up
-- +goose StatementBegin
alter table public.user
    add column if not exists nickname text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd