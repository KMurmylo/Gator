-- +goose Up
CREATE TABLE posts(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT con_posts_fk
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
);
CREATE INDEX posts_published_idx ON posts(published_at);
-- +goose Down
DROP TABLE posts;