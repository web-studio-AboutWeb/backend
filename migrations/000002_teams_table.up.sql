CREATE TABLE teams
(
    id          serial PRIMARY KEY NOT NULL,
    title       text               NOT NULL,
    description text               NOT NULL,
    image_id    text               NOT NULL,
    created_at  timestamptz        NOT NULL DEFAULT now(),
    updated_at  timestamptz        NOT NULL DEFAULT now(),
    disabled_at timestamptz
);