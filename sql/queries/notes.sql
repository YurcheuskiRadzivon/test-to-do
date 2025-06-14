-- name: CreateNote :one
WITH inserted_note AS (
    INSERT INTO notes (title, description, status, author_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id
),
usersnotes_insert AS (
    INSERT INTO usersnotes (user_id, note_id)
    SELECT $4, inserted_note.id FROM inserted_note
)
SELECT id FROM inserted_note;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1 AND author_id = $2;

-- name: UpdateNote :exec
UPDATE notes
SET title = $2, description = $3, status = $4
WHERE id = $1 AND author_id = $5;

-- name: GetNote :one
SELECT * FROM notes WHERE id = $1 AND author_id = $2;

-- name: GetNotes :many
SELECT * FROM notes
WHERE author_id = $1;