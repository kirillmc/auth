-- +goose Up
create table users
(
    id         serial primary key,
    username   text      not null unique,
    email      text      not null unique,
    password   text      not null,
    role       integer   not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);

create table roles_to_endpoints
(
    id       serial primary key,
    endpoint text    not null,
    role     integer not null
);

insert into roles_to_endpoints (endpoint, role)
values ('/user_v1.UserV1/Create', 2),
       ('/user_v1.UserV1/Get', 2),
       ('/user_v1.UserV1/Delete', 2),
       ('/user_v1.UserV1/Update', 2);
-- +goose Down
drop table users;

