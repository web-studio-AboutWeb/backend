CREATE TABLE teams
(
    id          serial PRIMARY KEY NOT NULL,
    title       text               NOT NULL,
    created_at  timestamptz        NOT NULL DEFAULT now(),
    updated_at  timestamptz,
    disabled_at timestamptz
);