CREATE TABLE projects
(
    id          serial PRIMARY KEY,
    title       text NOT NULL,
    description text NOT NULL,
    started_at  text NOT NULL DEFAULT now(),
    updated_at  text NOT NULL DEFAULT now(),
    cover_id    text,
    ended_at    text,
    link        text
);