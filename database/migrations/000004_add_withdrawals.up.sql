create table WITHDRAWALS
(
    WITHDRAWALS_NR           bigint primary key,
    USERS_ID                 bigint      not null references USERS (USERS_ID),
    WITHDRAWALS_SUM          bigint     not null,
    WITHDRAWALS_REQUESTED_AT timestamptz not null default current_timestamp
);