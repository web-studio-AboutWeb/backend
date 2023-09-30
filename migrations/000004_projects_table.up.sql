CREATE TABLE projects
(
    id          serial PRIMARY KEY,
    title       text        NOT NULL,
    description text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    team_id     int4 REFERENCES teams (id),
    image_id    text NOT NULL,
    started_at  timestamptz,
    ended_at    timestamptz,
    isactive    bool,
    link        text
);