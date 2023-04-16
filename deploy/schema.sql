create table if not exists projects
(
  id           smallserial primary key,
  title        text        not null,
  description  text        not null,
  cover_id     text,
  started_at   timestamptz not null,
  ended_at     timestamptz,
  link         text
);

create table if not exists project_documents
(
  id         bigserial not null,
  project_id smallint  not null references projects (id) on delete cascade
);

create type user_role as enum ('global admin', 'admin', 'moderator', 'user');

create table if not exists users
(
  id          smallserial   not null primary key,
  name        text          not null,
  surname     text          not null,
  login       text          not null,
  password    text          not null,
  role        user_role     not null default 'user',
  created_at  timestamptz   not null default now(),
  disabled_at timestamptz
);

create type staffer_position as enum ('frontend', 'backend', 'teamlead', 'manager', 'marketer');

create table if not exists staffers
(
  id          smallserial      not null primary key,
  user_id     smallint         not null references users (id) on delete cascade,
  project_id  smallint         not null references projects (id) on delete cascade,
  position    staffer_position not null
);

insert into users(name, surname, login, password, role)
values ('Денис', 'Камчатов', 'dkamchatov', 'something', 'admin')
on conflict do nothing;