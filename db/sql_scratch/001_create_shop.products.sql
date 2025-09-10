
CREATE SCHEMA public;
CREATE SCHEMA shop;
CREATE TABLE  shop.products (
id BIGSERIAL PRIMARY KEY ,--
sku TEXT NOT NULL UNIQUE,
name TEXT NOT NULL,
price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
stock_qty INTEGER NOT NULL DEFAULT 0 CHECK (stock_qty >= 0),
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
SELECT * FROM products;