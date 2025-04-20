-- name: GetLinksByUser :many
SELECT * FROM links WHERE user_id = $1;

-- name: GetLinkByShortCode :one
SELECT * FROM links WHERE short_code = $1;

-- name: CreateLink :exec
INSERT INTO links (
  full_url, short_code, user_id
) VALUES (
  $1, $2, $3
);

-- name: UpdateLink :exec
UPDATE links
SET short_code = $1, full_url = $2
WHERE id = $3;

-- name: UpdateClickCount :exec
UPDATE links
SET click_count = click_count + 1
WHERE id = $1;
