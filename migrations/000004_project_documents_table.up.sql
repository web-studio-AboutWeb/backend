CREATE TABLE IF NOT EXISTS project_documents
(
    id         BIGSERIAL PRIMARY KEY,
    project_id SMALLINT NOT NULL REFERENCES projects (id) ON DELETE CASCADE
);