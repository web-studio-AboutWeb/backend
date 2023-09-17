CREATE DATABASE ws WITH OWNER webstudio;

CREATE TABLE IF NOT EXISTS projects
(
    id          SMALLSERIAL PRIMARY KEY,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL,
    cover_id    TEXT,
    started_at  timestamptz NOT NULL,
    ended_at    timestamptz,
    link        TEXT
);

CREATE TABLE IF NOT EXISTS project_documents
(
    id         BIGSERIAL PRIMARY KEY,
    project_id SMALLINT NOT NULL REFERENCES projects (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users
(
    id          SMALLSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    surname     TEXT        NOT NULL,
    login       TEXT        NOT NULL,
    password    TEXT        NOT NULL,
    role        int2        NOT NULL,
    position    int2        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    disabled_at timestamptz
);

CREATE TABLE IF NOT EXISTS project_participants
(
    project_id SMALLINT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    user_id    SMALLINT NOT NULL REFERENCES users (id) ON DELETE CASCADE
);