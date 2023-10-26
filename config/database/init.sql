CREATE USER db_payments WITH CREATEDB NOSUPERUSER INHERIT PASSWORD 'password';
CREATE DATABASE db_payments OWNER db_payments;

\c db_payments db_payments
CREATE TABLE users(
    id VARCHAR primary key ,
    email VARCHAR (255) NOT NULL UNIQUE ,
    username VARCHAR (255) NOT NULL UNIQUE ,
    password VARCHAR (255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE customers (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR (255) references users(id),
    name VARCHAR (255) NOT NULL,
    balance BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE merchants (
    id VARCHAR PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    description TEXT,
    business_type VARCHAR (255),
    balance BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE transactions (
    id VARCHAR PRIMARY KEY,
    sender_customer_id VARCHAR REFERENCES customers (id),
    receiver_merchant_id VARCHAR REFERENCES merchants (id),
    amount BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);
