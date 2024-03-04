CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT UNSIGNED PRIMARY KEY,
    username VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_passwords (
    of_user_id BIGINT UNSIGNED PRIMARY KEY,
    hash VARCHAR(128) NOT NULL,
    FOREIGN KEY (of_user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS download_tasks (
    task_id BIGINT UNSIGNED PRIMARY KEY,
    of_user_id BIGINT UNSIGNED,
    download_type SMALLINT NOT NULL,
    url TEXT NOT NULL,
    download_status SMALLINT NOT NULL,
    metadata TEXT NOT NULL,
    FOREIGN KEY (of_user_id) REFERENCES users(user_id)
);