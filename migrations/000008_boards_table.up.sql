CREATE TABLE boards
(
    id         serial PRIMARY KEY,
    title      text NOT NULL,
    project_id int4 NOT NULL REFERENCES projects (id)
);