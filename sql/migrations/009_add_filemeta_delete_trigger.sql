-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION notify_filemeta_delete()
RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify('filemeta_events', OLD.uri);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER filemeta_delete_trigger
AFTER DELETE ON filemetas
FOR EACH ROW
EXECUTE FUNCTION notify_filemeta_delete();



