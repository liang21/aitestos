create type case_status_enum as enum ('unexecuted', 'pass', 'block', 'fail');

alter type case_status_enum owner to postgres;

create type case_type_enum as enum ('functionality', 'performance', 'api', 'ui', 'security');

alter type case_type_enum owner to postgres;

create type plan_status_enum as enum ('draft', 'active', 'completed', 'archived');

alter type plan_status_enum owner to postgres;

create type priority_enum as enum ('P0', 'P1', 'P2', 'P3');

alter type priority_enum owner to postgres;

create type result_status_enum as enum ('pass', 'fail', 'block', 'skip');

alter type result_status_enum owner to postgres;

create type user_role_enum as enum ('super_admin', 'admin', 'normal');

alter type user_role_enum owner to postgres;

create table project
(
    id          uuid                        default uuid_generate_v4() not null
        primary key,
    name        varchar(255)                                           not null
        unique,
    description text,
    config      jsonb                       default '{}'::jsonb,
    created_at  timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table project
    owner to postgres;

create table module
(
    id          uuid default uuid_generate_v4() not null
        primary key,
    project_id  uuid                            not null
        references project
            on delete cascade,
    name        varchar(255)                    not null,
    description text,
    unique (project_id, name)
);

alter table module
    owner to postgres;

create table users
(
    id         uuid                        default uuid_generate_v4()       not null
        primary key,
    username   varchar(32)                                                  not null,
    email      varchar(255)                                                 not null
        unique,
    password   varchar(255)                                                 not null,
    role       user_role_enum              default 'normal'::user_role_enum not null,
    created_at timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table users
    owner to postgres;

create table test_case
(
    id            uuid                        default uuid_generate_v4()              not null
        primary key,
    module_id     uuid                                                                not null
        references module
            on delete cascade,
    user_id       uuid                                                                not null
        references users,
    number        varchar(32)                                                         not null
        unique,
    title         varchar(255)                                                        not null,
    preconditions jsonb                       default '[]'::jsonb,
    steps         jsonb                       default '[]'::jsonb                     not null,
    expected      jsonb                       default '{}'::jsonb                     not null,
    ai_metadata   jsonb                       default '{}'::jsonb,
    case_type     case_type_enum              default 'functionality'::case_type_enum not null,
    priority      priority_enum               default 'P2'::priority_enum             not null,
    status        case_status_enum            default 'unexecuted'::case_status_enum  not null,
    created_at    timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at    timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table test_case
    owner to postgres;

create index idx_test_case_ai
    on test_case using gin (ai_metadata);

create index idx_test_case_steps
    on test_case using gin (steps);

create table test_plan
(
    id           uuid                        default uuid_generate_v4()        not null
        primary key,
    project_id   uuid                                                          not null
        references project
            on delete cascade,
    user_id      uuid                                                          not null
        references users,
    name         varchar(255)                                                  not null,
    status       plan_status_enum            default 'draft'::plan_status_enum not null,
    extra_config jsonb                       default '{}'::jsonb,
    created_at   timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at   timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table test_plan
    owner to postgres;

create table test_result
(
    id             uuid                        default uuid_generate_v4() not null
        primary key,
    case_id        uuid                                                   not null
        references test_case
            on delete cascade,
    plan_id        uuid                                                   not null
        references test_plan
            on delete cascade,
    executor_id    uuid                                                   not null
        references users,
    execute        result_status_enum                                     not null,
    result_details jsonb                       default '{}'::jsonb,
    executed_at    timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table test_result
    owner to postgres;

create table document
(
    id           uuid                        default uuid_generate_v4() not null
        primary key,
    project_id   uuid                                                   not null
        references project
            on delete cascade,
    name         varchar(255)                                           not null,
    type         varchar(32)                                            not null,
    url          text,
    content_text text,
    metadata     jsonb                       default '{}'::jsonb,
    created_at   timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at   timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table document
    owner to postgres;

create table document_chunk
(
    id          uuid                        default uuid_generate_v4() not null
        primary key,
    document_id uuid                                                   not null
        references document
            on delete cascade,
    chunk_index integer                                                not null,
    content     text                                                   not null,
    metadata    jsonb                       default '{}'::jsonb,
    created_at  timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table document_chunk
    owner to postgres;

create table generation_task
(
    id             uuid                        default uuid_generate_v4()           not null
        primary key,
    project_id     uuid                                                             not null
        references project
            on delete cascade,
    user_id        uuid                                                             not null
        references users,
    status         varchar(32)                 default 'pending'::character varying not null,
    prompt         text,
    result_summary jsonb                       default '{}'::jsonb,
    error_msg      text,
    created_at     timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at     timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table generation_task
    owner to postgres;

create table generated_case_draft
(
    id            uuid                        default uuid_generate_v4()           not null
        primary key,
    task_id       uuid                                                             not null
        references generation_task
            on delete cascade,
    module_id     uuid
                                                                                   references module
                                                                                       on delete set null,
    title         varchar(255)                                                     not null,
    preconditions jsonb                       default '[]'::jsonb,
    steps         jsonb                       default '[]'::jsonb                  not null,
    expected      jsonb                       default '{}'::jsonb                  not null,
    case_type     case_type_enum              default 'functionality'::case_type_enum,
    priority      priority_enum               default 'P2'::priority_enum,
    ai_metadata   jsonb                       default '{}'::jsonb,
    status        varchar(32)                 default 'pending'::character varying not null,
    feedback      text,
    created_at    timestamp(3) with time zone default CURRENT_TIMESTAMP,
    updated_at    timestamp(3) with time zone default CURRENT_TIMESTAMP
);

alter table generated_case_draft
    owner to postgres;

create function uuid_nil() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_nil() owner to postgres;

create function uuid_ns_dns() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_dns() owner to postgres;

create function uuid_ns_url() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_url() owner to postgres;

create function uuid_ns_oid() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_oid() owner to postgres;

create function uuid_ns_x500() returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_ns_x500() owner to postgres;

create function uuid_generate_v1() returns uuid
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v1() owner to postgres;

create function uuid_generate_v1mc() returns uuid
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v1mc() owner to postgres;

create function uuid_generate_v3(namespace uuid, name text) returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v3(uuid, text) owner to postgres;

create function uuid_generate_v4() returns uuid
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v4() owner to postgres;

create function uuid_generate_v5(namespace uuid, name text) returns uuid
    immutable
    strict
    parallel safe
    language c
as
$$
begin
-- missing source code
end;
$$;

alter function uuid_generate_v5(uuid, text) owner to postgres;

create function update_updated_at_column() returns trigger
    language plpgsql
as
$$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

alter function update_updated_at_column() owner to postgres;

