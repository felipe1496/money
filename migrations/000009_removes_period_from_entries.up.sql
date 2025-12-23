drop view if exists v_entries;

alter table entries drop column period;

create or replace view v_entries as
select
    e.id,
    e.transaction_id,
    t.name,
    t.description,
    e.amount,
    left(regexp_replace(e.reference_date::text, '[^0-9]', '', 'g'), 6) as period, -- Remove caracteres especiais e pega 6 dígitos
    t.user_id,
    t.category,
    sum(e.amount) over (partition by e.transaction_id) as total_amount,
    row_number() over (partition by e.transaction_id order by e.reference_date) as installment,
    count(*) over (partition by e.transaction_id) as total_installments,
    e.created_at,
    e.reference_date
from
    entries e
join transactions t on
    e.transaction_id = t.id;