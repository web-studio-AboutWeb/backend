CREATE TABLE documents
(
    id         serial PRIMARY KEY,
    filename   text        NOT NULL,
    file_id    text        NOT NULL,
    user_id    int4        NOT NULL REFERENCES users (id),
    mime       text        NOT NULL,
    size       int4        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);