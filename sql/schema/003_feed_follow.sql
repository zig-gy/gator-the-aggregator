-- +goose Up
CREATE TABLE feed_follow (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id),
    feed_id UUID REFERENCES feeds(id),
    CONSTRAINT composite_feed_user UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follow;