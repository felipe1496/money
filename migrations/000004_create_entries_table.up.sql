create table entries (
    id text primary key,
    transaction_id text not null references transactions(id),
    amount decimal(10,2) not null,
    period varchar(6) not null,
    created_at timestamptz not null default now()
);