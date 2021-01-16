alter table service
    add notifiers varchar[] default '{}'::varchar[] not null;
