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