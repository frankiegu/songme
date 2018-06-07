DROP INDEX IF EXISTS songs_confirmed_index;
DROP INDEX IF EXISTS songs_created_at_index;
ALTER TABLE songs DROP CONSTRAINT songs_user_id_fk;
-- DROP TABLE IF EXISTS songs;