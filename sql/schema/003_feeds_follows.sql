-- +goose Up
CREATE TABLE feed_follows(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT UN_feed_follows
    UNIQUE (user_id,feed_id),

    CONSTRAINT FK_ffollows_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

    CONSTRAINT FK_ffollows_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);
-- +goose Down
DROP TABLE feed_follows;