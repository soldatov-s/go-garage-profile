-- +goose Up
CREATE TABLE IF NOT EXISTS production.profile (
    id INTEGER PRIMARY KEY,
    user_first_name character varying(255) DEFAULT '',
    user_middle_name character varying(255) DEFAULT '',
    user_last_name character varying(255) DEFAULT '',
    user_position jsonb,
    user_company jsonb,
    user_private_key character varying(2048) DEFAULT '',
    user_public_key character varying(2048) DEFAULT '',
    meta jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL
);

-- +goose Down
DROP TABLE IF EXISTS production.profile;
