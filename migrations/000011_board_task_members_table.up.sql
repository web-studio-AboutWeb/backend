CREATE TABLE board_task_members
(
    task_id int4 NOT NULL REFERENCES board_tasks(id),
    user_id int4 NOT NULL REFERENCES users(id)
);