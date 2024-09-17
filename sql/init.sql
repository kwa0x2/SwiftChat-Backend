DO $$
BEGIN
CREATE TYPE public.friend_status AS ENUM
    ('friend', 'block_both', 'block_first_second', 'block_second_first');

CREATE TYPE public.read_status AS ENUM
    ('unread', 'readed');

CREATE TYPE public.request_status AS ENUM
    ('pending', 'rejected', 'accepted');

CREATE TYPE public.room_type AS ENUM
    ('private', 'group');

CREATE TABLE IF NOT EXISTS public."ROLE"
(
    role_name character varying(10) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "ROLE_pkey" PRIMARY KEY (role_name)
    );

CREATE TABLE IF NOT EXISTS public."USER"
(
    user_id character varying COLLATE pg_catalog."default" NOT NULL,
    user_email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    user_name character varying(10) COLLATE pg_catalog."default" NOT NULL,
    user_photo character varying COLLATE pg_catalog."default",
    user_role character varying(10) COLLATE pg_catalog."default" NOT NULL DEFAULT 'standard'::character varying,
    "createdAt" timestamp without time zone NOT NULL,
    "updatedAt" timestamp without time zone NOT NULL,
    "deletedAt" timestamp without time zone,
    CONSTRAINT "USER_pkey" PRIMARY KEY (user_id),
    CONSTRAINT "USER_user_email_key" UNIQUE (user_email),
    CONSTRAINT "USER_user_role_fkey" FOREIGN KEY (user_role)
    REFERENCES public."ROLE" (role_name) MATCH SIMPLE
                          ON UPDATE NO ACTION
                          ON DELETE NO ACTION
    );



CREATE TABLE IF NOT EXISTS public."ROOM"
(
    room_id uuid NOT NULL DEFAULT gen_random_uuid(),
    created_user_id character varying COLLATE pg_catalog."default" NOT NULL,
    message_count bigint,
    last_message text COLLATE pg_catalog."default",
    "updatedAt" timestamp without time zone NOT NULL,
    "createdAt" timestamp without time zone NOT NULL,
    "deletedAt" timestamp without time zone,
    room_type room_type NOT NULL DEFAULT 'private'::room_type,
    CONSTRAINT "ROOM_pkey" PRIMARY KEY (room_id),
    CONSTRAINT user_id FOREIGN KEY (created_user_id)
    REFERENCES public."USER" (user_id) MATCH SIMPLE
                          ON UPDATE NO ACTION
                          ON DELETE NO ACTION
    NOT VALID
    );

CREATE TABLE IF NOT EXISTS public."USER_ROOM"
(
    user_id character varying COLLATE pg_catalog."default" NOT NULL,
    "createdAt" timestamp without time zone NOT NULL,
    "updatedAt" timestamp without time zone NOT NULL,
    "deletedAt" timestamp without time zone,
    room_id uuid NOT NULL,
    CONSTRAINT "USER_ROOM_pkey" PRIMARY KEY (room_id, user_id),
    CONSTRAINT room_id FOREIGN KEY (room_id)
    REFERENCES public."ROOM" (room_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID,
    CONSTRAINT user_id FOREIGN KEY (user_id)
    REFERENCES public."USER" (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    );



CREATE TABLE IF NOT EXISTS public."MESSAGE"
(
    message text COLLATE pg_catalog."default" NOT NULL,
    sender_id character varying COLLATE pg_catalog."default" NOT NULL,
    message_status read_status NOT NULL DEFAULT 'unread'::read_status,
    "createdAt" timestamp without time zone NOT NULL,
    "updatedAt" timestamp without time zone NOT NULL,
    "deletedAt" timestamp without time zone,
    message_id uuid NOT NULL DEFAULT gen_random_uuid(),
    room_id uuid NOT NULL,
    CONSTRAINT "MESSAGE_pkey" PRIMARY KEY (message_id),
    CONSTRAINT room_id FOREIGN KEY (room_id)
    REFERENCES public."ROOM" (room_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID,
    CONSTRAINT sender_id FOREIGN KEY (sender_id)
    REFERENCES public."USER" (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID
    );

z
CREATE TABLE IF NOT EXISTS public."REQUEST"
(
    request_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    sender_mail character varying COLLATE pg_catalog."default" NOT NULL,
    receiver_mail character varying COLLATE pg_catalog."default" NOT NULL,
    request_status request_status NOT NULL DEFAULT 'pending'::request_status,
    "createdAt" timestamp without time zone NOT NULL,
    "deletedAt" timestamp without time zone,
    CONSTRAINT "REQUEST_pkey" PRIMARY KEY (request_id),
    CONSTRAINT request_unique UNIQUE NULLS NOT DISTINCT (sender_mail, receiver_mail, "deletedAt"),
    CONSTRAINT receiver_mail FOREIGN KEY (receiver_mail)
    REFERENCES public."USER" (user_email) MATCH SIMPLE
                          ON UPDATE NO ACTION
                          ON DELETE NO ACTION
    NOT VALID,
    CONSTRAINT sender_mail FOREIGN KEY (sender_mail)
    REFERENCES public."USER" (user_email) MATCH SIMPLE
                          ON UPDATE NO ACTION
                          ON DELETE NO ACTION
    NOT VALID
    );

INSERT INTO public."ROLE"(role_name)
VALUES
    ('high'),
    ('medium'),
    ('standard');

EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END $$;