create table "check"
(
    id              uuid        not null,
    service_id      uuid        not null
        constraint check_service_id_fk
            references service
            on delete cascade,
    latency_in_ms   int,
    ping_time_in_ms int,
    is_failure      bool,
    created_at      timestamptz not null
);

create unique index check_id_uindex
    on "check" (id);

create index check_service_id_created_at_index
    on "check" (service_id asc, created_at desc);

alter table "check"
    add constraint check_pk
        primary key (id);

