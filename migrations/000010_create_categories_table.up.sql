create table categories (
    id text primary key,
    user_id text not null,
    name text not null,
    color char(7) not null,
    created_at timestamptz not null default now()
);