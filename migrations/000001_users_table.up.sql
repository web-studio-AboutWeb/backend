BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id               SMALLSERIAL PRIMARY KEY,
    name             TEXT        NOT NULL,
    surname          TEXT        NOT NULL,
    login            TEXT        NOT NULL UNIQUE,
    encoded_password TEXT        NOT NULL,
    role             int2        NOT NULL,
    position         int2        NOT NULL,
    created_at       timestamptz NOT NULL DEFAULT now(),
    disabled_at      timestamptz
);
CREATE UNIQUE INDEX usr_login_idx ON users (login);

COMMIT;