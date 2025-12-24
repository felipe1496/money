CREATE OR REPLACE VIEW v_category_amount_per_period AS
SELECT 
    t.category_id as id,
    t.user_id,
    c.name,
    c.color,
    TO_CHAR(e.reference_date, 'YYYYMM') AS period,
    SUM(e.amount) AS total_amount
FROM 
    entries e
    JOIN transactions t ON e.transaction_id = t.id
    JOIN categories c ON t.category_id = c.id
GROUP BY 
    t.user_id,
    t.category_id,
    c.name,
    c.color,
    TO_CHAR(e.reference_date, 'YYYYMM'),
    TO_DATE(TO_CHAR(e.reference_date, 'YYYYMM'), 'YYYYMM');