-- +migrate Up
CREATE TABLE usersnotes (
    user_id INT NOT NULL,
    note_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id,note_id)
);