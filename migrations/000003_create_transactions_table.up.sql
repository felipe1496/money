create table transactions (
    id text primary key,
    user_id text not null references users(id),
    category varchar(100) not null,
    name varchar(100) not null,
    description varchar(400),
    created_at timestamptz not null default now()
);