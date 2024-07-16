DO $$
BEGIN
    CREATE TYPE public.friendship_status AS ENUM
        ('pending', 'accepted', 'rejected', 'email_sent', 'blocked');

    CREATE TYPE public.read_status AS ENUM
        ('unread', 'readed');

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
        CONSTRAINT "USER_user_role_fkey" FOREIGN KEY (user_role)
            REFERENCES public."ROLE" (role_name) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION
            NOT VALID
    );

    CREATE TABLE IF NOT EXISTS public."MESSAGE"
    (
        message_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
        message_content character varying(10000) COLLATE pg_catalog."default" NOT NULL,
        "createdAt" timestamp without time zone NOT NULL,
        "updatedAt" timestamp without time zone NOT NULL,
        "deletedAt" timestamp without time zone,
        message_sender_id character varying COLLATE pg_catalog."default" NOT NULL,
        message_receiver_id character varying COLLATE pg_catalog."default" NOT NULL,
        message_read_status read_status NOT NULL DEFAULT 'unread'::read_status,
        CONSTRAINT "MESSAGE_pkey" PRIMARY KEY (message_id),
        CONSTRAINT receiver_id FOREIGN KEY (message_receiver_id)
            REFERENCES public."USER" (user_id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION
            NOT VALID,
        CONSTRAINT sender_id FOREIGN KEY (message_sender_id)
            REFERENCES public."USER" (user_id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION
            NOT VALID
    );

    CREATE TABLE IF NOT EXISTS public."FRIENDSHIP"
    (
        sender_id character varying COLLATE pg_catalog."default" NOT NULL,
        receiver_id character varying COLLATE pg_catalog."default" NOT NULL,
        friendship_status friendship_status NOT NULL DEFAULT 'pending'::friendship_status,
        CONSTRAINT "FRIENDSHIP_pkey" PRIMARY KEY (sender_id, receiver_id),
        CONSTRAINT receiver_id FOREIGN KEY (receiver_id)
            REFERENCES public."USER" (user_id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION,
        CONSTRAINT sender_id FOREIGN KEY (sender_id)
            REFERENCES public."USER" (user_id) MATCH SIMPLE
            ON UPDATE NO ACTION
            ON DELETE NO ACTION
    );

    INSERT INTO public."ROLE"(role_name)
    VALUES
    ('high'),
    ('medium'),
    ('standard');

    COMMENT ON CONSTRAINT receiver_id ON public."MESSAGE"
        IS 'alıcı kişinin id değeri';

    COMMENT ON CONSTRAINT sender_id ON public."MESSAGE"
        IS 'mesajı yollayan kişinin id değeri';

EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END $$;
