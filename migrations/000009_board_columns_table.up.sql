CREATE TABLE board_columns
(
    id       serial PRIMARY KEY,
    board_id int4 NOT NULL REFERENCES boards (id),
    title    text NOT NULL
);