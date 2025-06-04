-- +goose Up
ALTER TABLE notes 
ADD COLUMN author_id INT NOT NULL DEFAULT 0;

ALTER TABLE notes
ADD CONSTRAINT fk_notes_users 
FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;

