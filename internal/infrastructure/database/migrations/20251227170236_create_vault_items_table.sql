-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vault_items(
    id SERIAL PRIMARY KEY, 
    name VARCHAR(50) NOT NULL,
    description TEXT,
    encrypted_username BYTEA NOT NULL,
    encrypted_password BYTEA NOT NULL,
    encrypted_url BYTEA,
    encrypted_note BYTEA,
    nonce BYTEA NOT NULL,
    creator_id INT REFERENCES accounts(id) ON DELETE CASCADE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (name, creator_id)
);

CREATE TABLE IF NOT EXISTS vault_items_groups (
    vault_item_id INT REFERENCES vault_items(id) ON DELETE CASCADE,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY(vault_item_id, group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vault_items_groups;
DROP TABLE IF EXISTS vault_items;
-- +goose StatementEnd
