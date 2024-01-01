CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email VARCHAR(320) NOT NULL UNIQUE,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE todo_lists
(
    id    SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE users_lists
(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    list_id INT REFERENCES todo_lists(id) ON DELETE CASCADE
);

CREATE TABLE todo_items
(
    id          SERIAL PRIMARY KEY,
    favorite    BOOLEAN               DEFAULT FALSE,
    description VARCHAR(255) NOT NULL,
    done        BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE lists_items
(
    id SERIAL PRIMARY KEY,
    item_id INT REFERENCES todo_items(id) ON DELETE CASCADE,
    list_id INT REFERENCES todo_lists(id) ON DELETE CASCADE
);