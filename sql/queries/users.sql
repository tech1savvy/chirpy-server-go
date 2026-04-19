-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
gen_random_uuid(),
  NOW(),
  NOW(),
  @email,
  @hashed_password
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = @email;

-- name: UpdateUser :one
UPDATE users
  SET
    email = $1,
    hashed_password = $2,
    updated_at = NOW()
  WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUser :one
UPDATE users
  SET is_chirpy_red = true, updated_at = NOW()
  WHERE id = $1
RETURNING *;
