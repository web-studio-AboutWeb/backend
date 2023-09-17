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