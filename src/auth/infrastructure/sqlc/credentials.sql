-- name: CreateCredentials :exec
INSERT INTO
    credentials (
        id,
        owner_id,
        email,
        password_hash,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6);

-- name: FindCredentialsByEmail :one
SELECT
    *
FROM
    credentials
WHERE
    email = $1;