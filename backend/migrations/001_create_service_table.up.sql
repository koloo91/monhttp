create table service
(
    id uuid not null,
    name varchar not null,
    type varchar not null,
    interval_in_seconds int not null,
    next_check_time timestamptz not null,
    endpoint varchar not null,
    http_method varchar,
    request_timeout_in_seconds int,
    http_headers varchar,
    http_body varchar,
    expected_http_response_body varchar,
    expected_http_status_code int,
    follow_redirects boolean,
    verify_ssl boolean,
    enable_notifications boolean default false not null,
    notify_after_number_of_failures int,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create unique index service_id_uindex
    on service (id);

create index service_next_check_time_index
    on service (next_check_time);

alter table service
    add constraint service_pk
        primary key (id);

