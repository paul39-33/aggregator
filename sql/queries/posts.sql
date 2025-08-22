-- name: CreatePost :one
INSERT INTO posts(title, url, description, published_at, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, users.id AS user_id, users.name AS user_name, feeds.url AS feed_url, feeds.name AS feed_name
FROM posts
INNER JOIN feeds
ON posts.feed_id = feeds.id
INNER JOIN users
ON feeds.user_id = users.id
WHERE users.id = $1
ORDER BY posts.created_at DESC
LIMIT $2;