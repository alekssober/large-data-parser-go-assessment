CREATE TABLE IF NOT EXISTS products (
id UUID PRIMARY KEY,
sku TEXT UNIQUE NOT NULL,
name TEXT NOT NULL,
category TEXT NOT NULL,
description TEXT NOT NULL,
brand_name TEXT NOT NULL,
stock_quantity INTEGER NOT NULL CHECK (stock_quantity >= 0),
manufacturer TEXT NOT NULL,
weight_grams INTEGER NOT NULL CHECK (weight_grams >= 0),
color TEXT NOT NULL,
price_cents INTEGER NOT NULL CHECK (price_cents >= 0),
currency CHAR(3) NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand_name);
CREATE INDEX IF NOT EXISTS idx_products_color ON products(color);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);