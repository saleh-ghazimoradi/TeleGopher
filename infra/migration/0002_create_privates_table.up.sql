CREATE TABLE IF NOT EXISTS privates (
    id BIGSERIAL PRIMARY KEY,
    user1_id BIGINT NOT NULL,
    user2_id BIGINT NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (user1_id, user2_id),
    CHECK (user1_id < user2_id),
    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_privates_user1_id ON privates(user1_id);
CREATE INDEX IF NOT EXISTS idx_privates_user2_id ON privates(user2_id);