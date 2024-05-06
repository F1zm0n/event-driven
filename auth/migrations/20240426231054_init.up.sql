CREATE TABLE customer (
  customer_id UUID PRIMARY KEY,
  email VARCHAR(256) UNIQUE NOT NULL,
  password BYTEA NOT NULL
);
