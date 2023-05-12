create table USER_SESSIONS
(
    USERS_ID                 bigint primary key references USERS (USERS_ID),
    USER_SESSIONS_SIG_KEY    bytea       not null,
    USER_SESSIONS_STARTED_AT timestamptz not null default current_timestamp
);