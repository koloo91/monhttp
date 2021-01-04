create table job
(
    id         uuid        not null,
    service_id uuid        not null
        constraint job_service_id_fk
            references service
            on delete cascade,
    execute_at timestamptz not null,
    created_at timestamptz not null,
    updated_at timestamptz not null
);

create index job_execute_at_index
    on job (execute_at);

create unique index job_id_uindex
    on job (id);

alter table job
    add constraint job_pk
        primary key (id);

