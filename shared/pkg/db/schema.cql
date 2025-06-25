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


CREATE TABLE IF NOT EXISTS orders (
    order_id bigint PRIMARY KEY,
    user_id bigint,
    status text,
    total_amount decimal,
    created_at timestamp,
    updated_at timestamp
);



CREATE TABLE IF NOT EXISTS order_items (
    order_item_id bigint PRIMARY KEY,
    order_id bigint,
    product_id bigint,
    quantity int,
    price decimal
);


CREATE TABLE IF NOT EXISTS inventory_reservations (
    reservation_id bigint PRIMARY KEY,
    order_id bigint,
    product_id bigint,
    quantity_reserved int,
    reserved_at timestamp,
    released boolean
);


CREATE TABLE IF NOT EXISTS payments (
    payment_id bigint PRIMARY KEY, 
    order_id bigint,
    amount decimal,
    status text,
    transaction_id text,
    paid_at timestamp
);



CREATE TABLE IF NOT EXISTS notifications (
    notification_id bigint PRIMARY KEY, 
    order_id bigint,
    type text,
    sent_at timestamp,
    status text
);