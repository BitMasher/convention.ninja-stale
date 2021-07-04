create table if not exists pronouns
(
    id            serial not null
        constraint pronouns_pk
            primary key,
    subject_label varchar(8),
    object_label  varchar(8)
);

create
    unique index if not exists pronouns_subject_object_uindex
    on pronouns (subject_label, object_label);

insert into pronouns (subject_label, object_label)
values ('he', 'him'),
       ('she', 'her'),
       ('it', 'it'),
       ('they', 'them'),
       ('one', 'one'),
       ('ey', 'em'),
       ('e', 'em'),
       ('xe', 'xem'),
       ('sie', 'hir'),
       ('ve', 'ver'),
       ('ze', 'zir'),
       ('ne', 'nem'),
       ('ze', 'zem'),
       ('fae', 'faer'),
       ('per', 'per')
on conflict (subject_label, object_label) do nothing;

create table if not exists users
(
    id           serial      not null
        constraint user_pk
            primary key,
    idp_id       varchar(64) not null,
    first_name   varchar(64),
    last_name    varchar(64),
    dob          date,
    display_name varchar(32),
    pronoun_id   integer
        constraint user_pronouns_id_fk
            references pronouns
);

create unique index if not exists user_idp_id_uindex
    on users (idp_id);

create table if not exists tos_agreements
(
    id      serial not null
        constraint tos_agreements_pk
            primary key,
    user_id integer,
    agreed  boolean,
    ts      timestamp
);

create index if not exists tos_agreements_user_id_index
    on tos_agreements (user_id);
