CREATE TABLE documents
(
    id         serial PRIMARY KEY,
    filename   text        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);