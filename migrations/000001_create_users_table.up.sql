create table users (
    id text primary key,
    email varchar(400) not null unique,
    name varchar(200) not null,
    avatar_url varchar(3000),
    created_at timestamptz not null default now()
);