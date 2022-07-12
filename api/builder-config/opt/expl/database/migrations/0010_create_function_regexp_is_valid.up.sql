CREATE OR REPLACE FUNCTION regexp_is_valid(text) RETURNS boolean AS
$$
BEGIN
    PERFORM regexp_match('', $1);
    RETURN true;
EXCEPTION
    WHEN invalid_regular_expression
        THEN RETURN false;
END;
$$ LANGUAGE plpgsql
