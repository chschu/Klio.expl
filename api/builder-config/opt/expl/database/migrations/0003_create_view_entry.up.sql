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
             ROW_NUMBER() OVER (PARTITION BY key_normalized, visible ORDER BY id)      head_index,
             ROW_NUMBER() OVER (PARTITION BY key_normalized, visible ORDER BY id DESC) tail_index,
             ROW_NUMBER() OVER (PARTITION BY key_normalized ORDER BY id)               permanent_index
      FROM entry_data) s
WHERE s.visible
