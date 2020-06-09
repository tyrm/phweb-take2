-- +migrate Up
CREATE TABLE "public"."users" (
    id serial NOT NULL UNIQUE,
    subject character varying NOT NULL UNIQUE,
    email character varying,
    name character varying,
    preferred_username character varying,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
)
;

-- +migrate Down
DROP TABLE "public"."users";