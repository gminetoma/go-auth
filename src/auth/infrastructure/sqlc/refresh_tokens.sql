-- name: CreateRefreshToken :exec
INSERT INTO
    refresh_tokens (
        id,
        owner_id,
        token_hash,
        expires_at,
        created_at,
        revoked_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6);

-- name: FindRefreshTokenByToken :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    token_hash = $1;

-- name: DeleteRefreshTokenByToken :exec
DELETE FROM refresh_tokens
WHERE
    token_hash = $1;

-- name: UpdateRefreshToken :execrows
UPDATE refresh_tokens
SET
    revoked_at = $1
WHERE
    id = $2;