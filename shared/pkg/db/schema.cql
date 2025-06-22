CREATE TABLE IF NOT EXISTS customers (
    id bigint PRIMARY KEY,
    username text,
    alias_name text,
    email text,
    created_at timestamp
);


CREATE TABLE products (
    id bigint PRIMARY KEY,
    name text,
    description text,
    price double,
    currency text,
    image_url text,
    stock int,
    created_at timestamp,
    updated_at timestamp
);
