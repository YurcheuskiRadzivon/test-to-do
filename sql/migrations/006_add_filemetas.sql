-- +migrate Up
CREATE TABLE filemetas (
    id SERIAL PRIMARY KEY,
    content_type VARCHAR(255) NOT NULL,
    owner_type VARCHAR(30) NOT NULL,
    owner_id INT NOT NULL,
    user_id INT NOT NULL,
    uri VARCHAR(350) NOT NULL
);

ALTER TABLE filemetas
ADD CONSTRAINT file_owner_type_check
CHECK (owner_type IN ('NOTE'))
NOT VALID;