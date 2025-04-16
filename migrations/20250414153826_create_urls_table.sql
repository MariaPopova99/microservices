-- +goose Up
create table urls(
    id serial primary key,
    short text not null unique,
    long text not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table urls;
