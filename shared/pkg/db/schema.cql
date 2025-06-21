CREATE TABLE IF NOT EXISTS customers (
    id bigint PRIMARY KEY,
    username text,
    alias_name text,
    email text,
    created_at timestamp
);
