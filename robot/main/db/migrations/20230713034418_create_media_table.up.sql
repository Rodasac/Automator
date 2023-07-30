CREATE TABLE IF NOT EXISTS media (
    id varchar(32) PRIMARY KEY,
    attributes jsonb,
    height double precision NOT NULL,
    width double precision NOT NULL,
    x double precision NOT NULL,
    y double precision NOT NULL,
    url varchar(255) NOT NULL,
    phash varchar(255) NOT NULL,
    filename varchar(255) NOT NULL,
    media_url varchar(255) NOT NULL,
    screenshot_url varchar(255) NOT NULL,
    resource_url varchar(255),
    task_id varchar(32) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);
