-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS groups(
    id SERIAL PRIMARY KEY, 
    name VARCHAR(50) NOT NULL,
    description TEXT,
    owner_id INT REFERENCES accounts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE accounts_groups (
    account_id INT REFERENCES accounts(id) ON DELETE CASCADE,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY(account_id, group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts_groups;
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
