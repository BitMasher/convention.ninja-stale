CREATE TABLE users
(
    id           varchar(16) NOT NULL,
    display_name varchar(32) NULL,
    "name"       varchar(64) NULL,
    dob          date        NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE user_oauth_provider
(
    id       varchar(256) NOT NULL,
    provider varchar(32)  NOT NULL,
    user_id  varchar(16)  NOT NULL,
    CONSTRAINT user_oauth_provider_un UNIQUE (user_id),
    CONSTRAINT user_oauth_provider_pk PRIMARY KEY (id, provider)
);

ALTER TABLE user_oauth_provider
    ADD CONSTRAINT user_oauth_provider_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;