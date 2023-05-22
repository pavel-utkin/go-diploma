alter table ORDERS
    add column ORDERS_PROCESSED_AT timestamptz not null default current_timestamp;