-- +migrate Up
ALTER TABLE filemetas
ADD CONSTRAINT fk_filemetas_user
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE;