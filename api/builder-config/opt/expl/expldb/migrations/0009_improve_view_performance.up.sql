CREATE OR REPLACE VIEW entry AS
SELECT s.id,             -- read-only
       s.key,
       s.key_normalized, -- read-only
       s.value,
       s.created_by,
       s.created_at,
       s.head_index,     -- read-only
       s.tail_index,     -- read-only
       s.permanent_index -- read-only
FROM (SELECT *,
             COUNT(*) FILTER (WHERE visible) OVER idx                                                    head_index,
             COUNT(*) FILTER (WHERE visible) OVER (idx ROWS BETWEEN CURRENT ROW AND UNBOUNDED FOLLOWING) tail_index,
             ROW_NUMBER() OVER idx                                                                       permanent_index
      FROM entry_data WINDOW idx AS (PARTITION BY key_normalized ORDER BY id)) s
WHERE s.visible
