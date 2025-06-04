-- name: CreateNote :exec
INSERT INTO notes (title, description, status)
VALUES ($1, $2, $3);

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;

-- name: UpdateNote :exec
UPDATE notes
SET title = $2, description = $3, status = $4
WHERE id = $1;

-- name: GetNote :one
SELECT * FROM notes WHERE id = $1;

-- name: GetNotes :many
SELECT * FROM notes;
