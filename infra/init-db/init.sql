CREATE TABLE IF NOT EXISTS sources
(
    id     SERIAL PRIMARY KEY,
    name   TEXT,
    type   TEXT,
    config JSONB
);

CREATE TABLE IF NOT EXISTS datasets
(
    id        SERIAL PRIMARY KEY,
    source_id INT NOT NULL REFERENCES sources (id),
    name      TEXT,
    config    JSONB
);

CREATE TABLE IF NOT EXISTS charts
(
    id         SERIAL PRIMARY KEY,
    dataset_id INT REFERENCES datasets (id),
    name       TEXT,
    type       TEXT,
    config     JSONB
);

CREATE TABLE IF NOT EXISTS dashboards
(
    id   SERIAL PRIMARY KEY,
    name TEXT,
    grid JSONB
);

-- sample schema
CREATE SCHEMA IF NOT EXISTS samples AUTHORIZATION admin;

CREATE TABLE IF NOT EXISTS samples.apple_stock
(
    time   DATE,
    open   DECIMAL,
    high   DECIMAL,
    low    DECIMAL,
    close  DECIMAL,
    volume DECIMAL
);

COPY samples.apple_stock
    FROM '/docker-entrypoint-initdb.d/apple_stock_price_2018_2024.csv'
    DELIMITER ','
    CSV HEADER;