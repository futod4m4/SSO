CREATE TABLE followers (
                           id SERIAL PRIMARY KEY,
                           follower_id INTEGER NOT NULL REFERENCES users(id),
                           following_id INTEGER NOT NULL REFERENCES users(id),
                           UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_follower_id ON followers(follower_id);
CREATE INDEX idx_following_id ON followers(following_id);