CREATE TABLE IF NOT EXISTS "USER" (
    user_id character varying NOT NULL,
    user_email character varying(50) NOT NULL,
    user_name character varying(10) NOT NULL,
    user_photo character varying,
    "createdAt" timestamp without time zone NOT NULL,
    "updatedAt" timestamp without time zone NOT NULL,
    "deletedAt" timestamp without time zone
);