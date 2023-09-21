CREATE TABLE team_members
(
    user_id    int4        NOT NULL REFERENCES users (id),
    team_id    int4        NOT NULL REFERENCES teams (id),
    role       int2        NOT NULL,
    position   int2        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);