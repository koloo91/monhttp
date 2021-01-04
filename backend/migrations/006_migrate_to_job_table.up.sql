INSERT INTO job (id, service_id, execute_at, created_at, updated_at)
SELECT gen_random_uuid(), id, next_check_time, now(), now()
FROM service;

drop index service_next_check_time_index;

alter table service drop column next_check_time;
