-- +goose Up
CREATE TABLE posts (
    id UUID NOT NULL PRIMARY KEY,
    feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url VARCHAR(256) NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP
);

-- +goose Down
DROP TABLE posts;
