-- name: CreateFileMeta :exec
INSERT INTO filemetas (
    content_type,
    owner_type,
    owner_id,
    user_id,
    uri
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);

-- name: DeleteFileMetaByID :exec
DELETE FROM filemetas
WHERE id = $1;

-- name: DeleteFileMetaByNoteID :exec
DELETE FROM filemetas
WHERE owner_type = $1 AND owner_id = $2;

-- name: GetFileMetaURI :one
SELECT uri
FROM filemetas
WHERE id = $1;

-- name: GetFileMetaByID :one
SELECT content_type, owner_type, owner_id, user_id,uri
FROM filemetas
WHERE id = $1;

-- name: GetFileMetaIDByID :many
SELECT id
FROM filemetas
WHERE owner_type = $1 AND owner_id = $2;

-- name: GetFileMetasIDByUserID :many
SELECT id
FROM filemetas
WHERE user_id = $1;

-- name: GetFileMetas :many
SELECT id, content_type, owner_type, owner_id,user_id, uri
FROM filemetas;

-- name: FileMetasExistsByIDAndUserID :one
SELECT EXISTS(SELECT 1 FROM filemetas WHERE id = $1 AND user_id = $2);