alter table entries add column reference_date date;

update entries 
set reference_date = to_date(period || '01', 'YYYYMMDD')
where reference_date is null;

alter table entries alter column reference_date set not null;