CREATE OR REPLACE TRIGGER entry_delete
    INSTEAD OF DELETE ON entry
    FOR EACH ROW
    EXECUTE FUNCTION entry_delete()
