create table ORDERS
(
    ORDERS_NR          bigint primary key,
    USERS_ID           bigint      not null references USERS (USERS_ID),
    ORDERS_STATUS text not null,
    ORDERS_UPLOADED_AT timestamptz not null default current_timestamp
);