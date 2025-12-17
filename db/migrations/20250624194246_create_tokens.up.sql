CREATE TABLE IF NOT EXISTS tokens
(
    jti        UUID PRIMARY KEY,
    token      TEXT      NOT NULL,
    user_id    BIGINT    NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tokens_player_id ON tokens (user_id);

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE
    ON tokens
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();