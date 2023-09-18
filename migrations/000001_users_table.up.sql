CREATE TABLE IF NOT EXISTS users
(
    id               SMALLSERIAL PRIMARY KEY,
    name             TEXT        NOT NULL,
    surname          TEXT        NOT NULL,
    email            TEXT        NOT NULL UNIQUE,
    username         TEXT        NOT NULL UNIQUE,
    encoded_password TEXT        NOT NULL,
    salt             TEXT        NOT NULL,
    role             int2        NOT NULL,
    position         int2        NOT NULL,
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz,
    disabled_at      timestamptz
);