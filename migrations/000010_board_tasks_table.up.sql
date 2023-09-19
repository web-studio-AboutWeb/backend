CREATE TABLE board_tasks
(
    id          serial PRIMARY KEY,
    board_id    int4    NOT NULL REFERENCES boards (id),
    column_id   int4    NOT NULL REFERENCES board_columns (id),
    title       text    NOT NULL,
    description text    NOT NULL,
    status      int2 NOT NULL
);