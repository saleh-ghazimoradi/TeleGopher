CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    from_id BIGINT NOT NULL,
    private_id BIGINT,
    message_type TEXT NOT NULL,
    content TEXT NOT NULL,
    delivered BOOLEAN NOT NULL DEFAULT false,
    read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (from_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (private_id) REFERENCES privates(id) ON DELETE CASCADE,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_messages_private_id ON messages(private_id);
CREATE INDEX IF NOT EXISTS idx_messages_from_id ON messages(from_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);