DO
$$
BEGIN

    CREATE TABLE IF NOT EXISTS "USER" (
        user_id character varying NOT NULL,
        user_email character varying(50) NOT NULL,
        user_name character varying(10) NOT NULL,
        user_photo character varying,
        user_role character varying(5) NOT NULL DEFAULT 'guest',
        "createdAt" timestamp without time zone NOT NULL,
        "updatedAt" timestamp without time zone NOT NULL,
        "deletedAt" timestamp without time zone
    );

    CREATE TABLE IF NOT EXISTS "ROLE" (
        role_name character varying(5) PRIMARY KEY
    );

    ALTER TABLE "USER"
    ADD FOREIGN KEY (user_role) REFERENCES "ROLE" (role_name);

    INSERT INTO "ROLE"(role_name)
    VALUES
    ('admin'),
    ('user'),
    ('guest');

EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END;
$$;