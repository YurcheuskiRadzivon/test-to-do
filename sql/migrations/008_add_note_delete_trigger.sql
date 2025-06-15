-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION delete_note_filemetas()
RETURNS TRIGGER AS $func$
BEGIN
    DELETE FROM filemetas
    WHERE owner_type = 'NOTE'
    AND owner_id = OLD.id;
    RETURN OLD;
END;
$func$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER trg_delete_note_filemetas
AFTER DELETE ON notes
FOR EACH ROW
EXECUTE FUNCTION delete_note_filemetas();

