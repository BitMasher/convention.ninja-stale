CREATE TABLE users
(
    id           varchar(27)  NOT NULL,
    display_name varchar(32)  NULL,
    first_name   varchar(128) NULL,
    last_name    varchar(128) NULL,
    dob          date         NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE TABLE user_oauth_providers
(
    id       varchar(256) NOT NULL,
    provider varchar(32)  NOT NULL,
    user_id  varchar(27)  NOT NULL,
    CONSTRAINT user_oauth_provider_un UNIQUE (user_id),
    CONSTRAINT user_oauth_provider_pk PRIMARY KEY (id, provider)
);

ALTER TABLE user_oauth_providers
    ADD CONSTRAINT user_oauth_provider_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;