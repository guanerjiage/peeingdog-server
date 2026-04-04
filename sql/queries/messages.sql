-- name: CreateMessage :one
INSERT INTO messages (user_id, text, latitude, longitude)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, text, latitude, longitude, created_at, expires_at;

-- name: GetMessageByID :one
SELECT id, user_id, text, latitude, longitude, created_at, expires_at
FROM messages
WHERE id = $1 AND expires_at > CURRENT_TIMESTAMP;

-- name: GetUserMessages :many
SELECT id, user_id, text, latitude, longitude, created_at, expires_at
FROM messages
WHERE user_id = $1 AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;

-- name: GetNearbyMessages :many
SELECT id, user_id, text, latitude, longitude, created_at, expires_at
FROM messages
WHERE latitude BETWEEN $1 AND $2
  AND longitude BETWEEN $3 AND $4
  AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;

-- name: GetExpiredMessages :many
SELECT id, user_id, text, latitude, longitude, created_at, expires_at
FROM messages
WHERE expires_at <= CURRENT_TIMESTAMP
ORDER BY expires_at ASC;

-- name: ArchiveExpiredMessages :exec
SELECT archive_expired_messages();

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1;
