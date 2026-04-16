-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
gen_random_uuid(),
  NOW(),
  NOW(),
  @email,
  @hashed_password
)
RETURNING id, created_at, updated_at, email;

