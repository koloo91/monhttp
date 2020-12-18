create table failure
(
    id uuid not null,
    service_id uuid not null
        constraint failure_service_id_fk
            references service
            on delete cascade,
    reason varchar not null,
    created_at timestamptz not null
);

create unique index failure_id_uindex
    on failure (id);

alter table failure
    add constraint failure_pk
        primary key (id);

create index failure_id_created_at_index
    on failure (id asc, created_at desc);

