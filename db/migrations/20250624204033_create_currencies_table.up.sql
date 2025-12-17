CREATE TABLE currencies
(
    id            BIGSERIAL PRIMARY KEY,
    code          VARCHAR(10)   NOT NULL DEFAULT 'USD',
    name          VARCHAR(100)  NOT NULL DEFAULT 'US Dollar',
    symbol        VARCHAR(10)   NOT NULL DEFAULT '$',
    precision     INT           NOT NULL DEFAULT 2,
    is_crypto     BOOLEAN       NOT NULL DEFAULT FALSE,
    is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP     NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_currencies_unique ON currencies (code);

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE
    ON currencies
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();