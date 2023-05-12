create table USERS
(
    USERS_ID            bigserial primary key,
    USERS_LOGIN         text unique not null,
    USERS_PASSWORD_HASH bytea not null
);