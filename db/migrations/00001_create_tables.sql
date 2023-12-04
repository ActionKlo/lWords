-- +goose Up
CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    eng text UNIQUE NOT NULL,
    rus text NOT NULL,
    learn_at timestamp DEFAULT (now())
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username text,
    password text,
    created_at timestamp
);
--
CREATE TABLE statistics (
    id SERIAL PRIMARY KEY,
    user_id integer,
    learned_words integer,
    not_learned_words integer,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE results (
    id SERIAL PRIMARY KEY,
    user_id integer,
    word_id integer,
    grade integer DEFAULT 2,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (word_id) REFERENCES words (id)
);

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
DROP TABLE statistics;

DROP TABLE results;

DROP TABLE words;

DROP TABLE users;

-- +goose StatementBegin
-- +goose StatementEnd
