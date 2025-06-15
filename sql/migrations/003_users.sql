-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(319) UNIQUE NOT NULL
);

INSERT INTO users (id, username, password, email)
VALUES (0,'admin', '$2a$10$UGGVQQm2CSuqsJnXoylLkOZGNyRp8nIGp9WU0uI4YLTx0jfvzr5ES', 'admin@example.com');
