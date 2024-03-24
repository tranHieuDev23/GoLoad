-- +migrate Up
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    account_name VARCHAR(256) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (account_name)
);

CREATE TABLE IF NOT EXISTS account_passwords (
    of_account_id BIGINT UNSIGNED AUTO_INCREMENT,
    hash VARCHAR(128) NOT NULL,
    PRIMARY KEY (of_account_id),
    FOREIGN KEY (of_account_id) REFERENCES accounts(id)
);

CREATE TABLE IF NOT EXISTS token_public_keys (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    public_key TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS download_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT,
    of_account_id BIGINT UNSIGNED,
    download_type SMALLINT NOT NULL,
    url TEXT NOT NULL,
    download_status SMALLINT NOT NULL,
    metadata TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (of_account_id) REFERENCES accounts(id)
);

-- +migrate Down
DROP TABLE IF EXISTS download_tasks;

DROP TABLE IF EXISTS token_public_keys;

DROP TABLE IF EXISTS account_passwords;

DROP TABLE IF EXISTS accounts;