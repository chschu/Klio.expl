CREATE OR REPLACE FUNCTION entry_delete() RETURNS trigger AS
$$
BEGIN
    UPDATE entry_data SET visible = false WHERE id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql
