-- +goose Up
alter table users
    rename column name to username;

alter table if exists users
    ADD UNIQUE (username, email);



create table if not exists roles_to_endpoints
(
    id       serial primary key,
    endpoint text    not null,
    role     integer not null
);

insert into roles_to_endpoints (endpoint, role)
values ('/chat_v1.ChatV1/Create', 2),
       ('/user_v1.UserV1/Get', 2),
       ('/user_v1.UserV1/Delete', 2),
       ('/user_v1.UserV1/Update', 2);
-- +goose Down
drop table if exists roles_to_endpoints;

