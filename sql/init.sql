DO $$
BEGIN
CREATE TYPE public.friend_status AS ENUM
    ('friend', 'blocked', 'unfriend');

CREATE TYPE public.read_status AS ENUM
    ('unread', 'readed');

CREATE TYPE public.request_status AS ENUM
    ('pending', 'rejected', 'accepted');

CREATE TYPE public.message_type AS ENUM
    ('text', 'photo', 'file');


CREATE TABLE IF NOT EXISTS public."USER"
(
    user_id character varying COLLATE pg_catalog."default" NOT NULL,
    user_email character varying(50) COLLATE pg_catalog."default" NOT NULL,
    user_name character varying(10) COLLATE pg_catalog."default" NOT NULL,
    user_photo character varying COLLATE pg_catalog."default",
    "createdAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "USER_pkey" PRIMARY KEY (user_id),
    CONSTRAINT "USER_user_email_key" UNIQUE (user_email)
    );

CREATE TABLE IF NOT EXISTS public."ROOM"
(
    room_id uuid NOT NULL DEFAULT gen_random_uuid(),
    created_user_id character varying COLLATE pg_catalog."default" NOT NULL,
    last_message_id uuid,
    "updatedAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "createdAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "ROOM_pkey" PRIMARY KEY (room_id),
    CONSTRAINT user_id FOREIGN KEY (created_user_id)
    REFERENCES public."USER" (user_id) MATCH SIMPLE
                          ON UPDATE NO ACTION
                          ON DELETE NO ACTION
    );


CREATE TABLE IF NOT EXISTS public."USER_ROOM"
(
    user_id character varying COLLATE pg_catalog."default" NOT NULL,
    room_id uuid NOT NULL,
    "createdAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "USER_ROOM_pkey" PRIMARY KEY (room_id, user_id),
    CONSTRAINT room_id FOREIGN KEY (room_id)
    REFERENCES public."ROOM" (room_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
    CONSTRAINT user_id FOREIGN KEY (user_id)
    REFERENCES public."USER" (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS public."MESSAGE"
(
    message_id uuid NOT NULL DEFAULT gen_random_uuid(),
    message_content text COLLATE pg_catalog."default" NOT NULL,
    sender_id character varying COLLATE pg_catalog."default" NOT NULL,
    room_id uuid NOT NULL,
    message_read_status read_status NOT NULL DEFAULT 'unread'::read_status,
    message_type message_type NOT NULL DEFAULT 'text'::message_type,
    message_starred boolean NOT NULL DEFAULT false,
    "createdAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "MESSAGE_pkey" PRIMARY KEY (message_id),
    CONSTRAINT room_id FOREIGN KEY (room_id)
    REFERENCES public."ROOM" (room_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
    CONSTRAINT sender_id FOREIGN KEY (sender_id)
    REFERENCES public."USER" (user_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS public."FRIEND"
(
    friend_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_email character varying COLLATE pg_catalog."default" NOT NULL,
    user_email2 character varying COLLATE pg_catalog."default" NOT NULL,
    friend_status friend_status NOT NULL DEFAULT 'friend'::friend_status,
    "createdAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "FRIEND_pkey" PRIMARY KEY (friend_id),
    CONSTRAINT friend_unique UNIQUE NULLS NOT DISTINCT (user_email, user_email2, "deletedAt"),
    CONSTRAINT user_email FOREIGN KEY (user_email)
    REFERENCES public."USER" (user_email) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
    CONSTRAINT user_email2 FOREIGN KEY (user_email2)
    REFERENCES public."USER" (user_email) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    );

CREATE TABLE IF NOT EXISTS public."REQUEST"
(
    request_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    sender_email character varying COLLATE pg_catalog."default" NOT NULL,
    receiver_email character varying COLLATE pg_catalog."default" NOT NULL,
    request_status request_status NOT NULL DEFAULT 'pending'::request_status,
    "createdAt" timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deletedAt" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "REQUEST_pkey" PRIMARY KEY (request_id),
    CONSTRAINT request_unique UNIQUE NULLS NOT DISTINCT (sender_email, receiver_email, "deletedAt"),
    CONSTRAINT sender_email FOREIGN KEY (sender_email)
    REFERENCES public."USER" (user_email) MATCH SIMPLE
                          ON UPDATE NO ACTION
                          ON DELETE NO ACTION
    );


EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END $$;