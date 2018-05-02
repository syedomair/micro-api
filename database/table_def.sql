
drop table public.clients CASCADE;
CREATE TABLE public.clients (
        id                  uuid PRIMARY KEY,
        name                varchar(100) NULL,
        api_key             varchar(100) NULL,
        secret              varchar(100) NULL,
        status              smallint NULL,
        created_at          timestamp,
        updated_at          timestamp
);

INSERT INTO public.clients (id, name, api_key, secret, status) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'MY Client', 'the$network#api*key', 'the$network#api*secret', 1 );

drop table public.batch_tasks CASCADE;
CREATE TABLE public.batch_tasks (
        id                  uuid PRIMARY KEY,
        client_id           uuid references public.clients(id),
        api_name            varchar(100) NULL,
        status              smallint NULL,
        created_at          timestamp,
        completed_at        timestamp NULL
);

drop table public.roles CASCADE;
CREATE TABLE public.roles (
        id                  uuid PRIMARY KEY,
        client_id           uuid references public.clients(id),
        title               varchar(100) NULL,
        role_type           smallint NULL,
        created_at          timestamp,
        updated_at          timestamp
);

drop table public.users CASCADE;
CREATE TABLE public.users (
        id                  uuid PRIMARY KEY,
        client_id           uuid references public.clients(id),
        first_name          varchar(100) NOT NULL,
        last_name           varchar(100) NOT NULL,
        email               varchar(100) NOT NULL,
        password            varchar(100) NOT NULL,
        is_admin            smallint NOT NULL,
        created_at          timestamp,
        updated_at          timestamp
);
