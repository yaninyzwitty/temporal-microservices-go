-- Customer Service CQL Schema
-- Based on customers.proto

-- Create keyspace for the customer service
CREATE KEYSPACE IF NOT EXISTS customer_service
WITH replication = {
    'class': 'SimpleStrategy',
    'replication_factor': 3
};

-- Use the keyspace
USE customer_service;

-- Main customers table
-- Primary key: id (for efficient lookups by customer ID)
CREATE TABLE IF NOT EXISTS customers (
    id bigint PRIMARY KEY,
    username text,
    alias_name text,
    email text,
    created_at timestamp,
    updated_at timestamp
);

-- Secondary index for username lookups
CREATE INDEX IF NOT EXISTS ON customers (username);

-- Secondary index for email lookups
CREATE INDEX IF NOT EXISTS ON customers (email);

-- Table for username to customer ID mapping (for unique username constraints)
CREATE TABLE IF NOT EXISTS customers_by_username (
    username text PRIMARY KEY,
    customer_id bigint
);

-- Table for email to customer ID mapping (for unique email constraints)
CREATE TABLE IF NOT EXISTS customers_by_email (
    email text PRIMARY KEY,
    customer_id bigint
);

-- Table for customer lookups by creation date (for analytics/reporting)
CREATE TABLE IF NOT EXISTS customers_by_created_at (
    created_at timestamp,
    customer_id bigint,
    username text,
    alias_name text,
    email text,
    updated_at timestamp,
    PRIMARY KEY (created_at, customer_id)
) WITH CLUSTERING ORDER BY (customer_id ASC);

-- Table for customer lookups by alias name (for search functionality)
CREATE TABLE IF NOT EXISTS customers_by_alias (
    alias_name text,
    customer_id bigint,
    username text,
    email text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (alias_name, customer_id)
) WITH CLUSTERING ORDER BY (customer_id ASC);

-- Comments for documentation
COMMENT ON TABLE customers IS 'Main customers table with customer_id as primary key';
COMMENT ON TABLE customers_by_username IS 'Lookup table for username to customer_id mapping';
COMMENT ON TABLE customers_by_email IS 'Lookup table for email to customer_id mapping';
COMMENT ON TABLE customers_by_created_at IS 'Time-series table for customer analytics and reporting';
COMMENT ON TABLE customers_by_alias IS 'Search table for customer lookups by alias name';

-- Sample data insertion (optional - for testing)
-- INSERT INTO customers (id, username, alias_name, email, created_at, updated_at)
-- VALUES (1, 'john_doe', 'John Doe', 'john.doe@example.com', toTimestamp(now()), toTimestamp(now()));

-- INSERT INTO customers_by_username (username, customer_id)
-- VALUES ('john_doe', 1);

-- INSERT INTO customers_by_email (email, customer_id)
-- VALUES ('john.doe@example.com', 1); 