-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT UNSIGNED PRIMARY KEY,
    account_name VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS account_passwords (
    of_account_id BIGINT UNSIGNED PRIMARY KEY,
    hash VARCHAR(128) NOT NULL,
    FOREIGN KEY (of_account_id) REFERENCES accounts(id)
);

CREATE TABLE IF NOT EXISTS token_public_keys (
    id BIGINT UNSIGNED PRIMARY KEY,
    public_key VARBINARY(4096) NOT NULL
);

CREATE TABLE IF NOT EXISTS download_tasks (
    task_id BIGINT UNSIGNED PRIMARY KEY,
    of_account_id BIGINT UNSIGNED,
    download_type SMALLINT NOT NULL,
    url TEXT NOT NULL,
    download_status SMALLINT NOT NULL,
    metadata TEXT NOT NULL,
    FOREIGN KEY (of_account_id) REFERENCES accounts(id)
);
-- +migrate StatementEnd

-- +migrate Down

-- +migrate StatementBegin
DROP TABLE IF EXISTS download_tasks;
DROP TABLE IF EXISTS token_public_keys;
DROP TABLE IF EXISTS account_passwords;
DROP TABLE IF EXISTS accounts;
-- +migrate StatementEnd
