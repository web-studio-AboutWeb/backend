CREATE TABLE IF NOT EXISTS project_categories
(
    id   smallserial PRIMARY KEY,
    name text NOT NULL UNIQUE
);