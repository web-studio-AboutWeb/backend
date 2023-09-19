CREATE TABLE project_documents
(
    project_id  int4 NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    document_id int4 NOT NULL REFERENCES documents (id) ON DELETE CASCADE
);