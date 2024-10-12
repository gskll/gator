-- +goose Up
CREATE TABLE feeds (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(32) NOT NULL,
    url VARCHAR(64) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE feeds;
