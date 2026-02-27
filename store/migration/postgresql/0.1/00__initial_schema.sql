CREATE TABLE public.users (
	id bigserial NOT NULL,
	username varchar(50) NOT NULL,
	nickname varchar(50) NULL,
	"password" varchar(255) NOT NULL,
	phone varchar(20) NULL,
	email varchar(100) NULL,
	"role" varchar(20) DEFAULT 'user'::character varying NULL,
	password_expires timestamptz NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);
CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);
CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);