-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    user_id UUID NOT NULL ,
    feed_id UUID NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
