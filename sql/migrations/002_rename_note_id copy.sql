-- +migrate Up
ALTER TABLE notes RENAME COLUMN note_id to id;