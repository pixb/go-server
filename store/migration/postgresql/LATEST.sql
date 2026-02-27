-- system_setting
CREATE TABLE public.system_setting (
  name varchar(255) NOT NULL,
  value text NOT NULL,
  description text NOT NULL DEFAULT '',
  CONSTRAINT system_setting_pkey PRIMARY KEY (name)
);

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
CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);
CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);

-- refresh_tokens table for PostgreSQL

CREATE TABLE public.refresh_tokens (
    id bigserial NOT NULL,
    user_id bigint NOT NULL,
    token text NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked boolean DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz NULL,
    CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id),
    CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE UNIQUE INDEX idx_refresh_tokens_token ON public.refresh_tokens USING btree (token);
CREATE INDEX idx_refresh_tokens_deleted_at ON public.refresh_tokens USING btree (deleted_at);
CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);
