CREATE TABLE IF NOT EXISTS public."user"
(
    user_id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default" NOT NULL,
    create_at timestamp with time zone DEFAULT now(),
    update_at timestamp with time zone DEFAULT now(),
    CONSTRAINT user_pk PRIMARY KEY (user_id),
    CONSTRAINT user_name_unique UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS public.task
(
    task_id uuid NOT NULL DEFAULT gen_random_uuid(),
    title character varying COLLATE pg_catalog."default" NOT NULL,
    description character varying COLLATE pg_catalog."default" NOT NULL,
    status character varying COLLATE pg_catalog."default" NOT NULL DEFAULT 'new'::character varying,
    create_at timestamp with time zone DEFAULT now(),
    update_at timestamp with time zone DEFAULT now(),
    CONSTRAINT task_pk PRIMARY KEY (task_id),
    CONSTRAINT tasks_status_check CHECK (status::text = ANY (ARRAY['new'::character varying, 'in_progress'::character varying, 'done'::character varying]::text[]))
);

CREATE TABLE IF NOT EXISTS public.user_task
(
    user_task_id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    task_id uuid NOT NULL,
    create_at timestamp with time zone DEFAULT now(),
    update_at timestamp with time zone DEFAULT now(),
    CONSTRAINT user_task_pk PRIMARY KEY (user_task_id),
    CONSTRAINT user_task_user_id_task_id_unique UNIQUE (user_id, task_id),
    CONSTRAINT user_task_task_fk FOREIGN KEY (task_id)
        REFERENCES public.task (task_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT user_task_user_fk FOREIGN KEY (user_id)
        REFERENCES public."user" (user_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);
CREATE INDEX IF NOT EXISTS user_task_task_id_btree
    ON public.user_task USING btree
    (task_id ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
CREATE INDEX IF NOT EXISTS user_task_user_id_btree
    ON public.user_task USING btree
    (user_id ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;