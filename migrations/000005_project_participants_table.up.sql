CREATE TABLE project_participants
(
    project_id int4        NOT NULL REFERENCES projects (id),
    user_id    int4        NOT NULL REFERENCES users (id),
    role       int2        NOT NULL,
    position   int2        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);