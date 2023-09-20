CREATE TABLE documents
(
    id         serial PRIMARY KEY,
    filename   text        NOT NULL,
    mime       text        NOT NULL,
    size       int4        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);