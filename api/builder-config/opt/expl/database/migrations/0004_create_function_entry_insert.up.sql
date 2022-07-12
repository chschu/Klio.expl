CREATE OR REPLACE FUNCTION entry_insert() RETURNS trigger AS
$$
DECLARE
    inserted_id INTEGER;
BEGIN
    INSERT INTO entry_data(key, key_normalized, value, created_by, created_at, visible)
    VALUES (NEW.key, NORMALIZE(LOWER(NEW.key), NFC), NEW.value, NEW.created_by, NEW.created_at, true)
    RETURNING id INTO inserted_id;

    SELECT * FROM entry WHERE id = inserted_id INTO NEW;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql
