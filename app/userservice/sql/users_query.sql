-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUser :many
SELECT * FROM users limit ?;
-- ORDER BY name;

-- name: CreateUser :execresult
INSERT INTO users (
  username, password, telnumber, email, salt, gender, status
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;