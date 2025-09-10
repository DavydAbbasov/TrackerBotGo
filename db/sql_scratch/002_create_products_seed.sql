INSERT INTO shop.products (sku,name,price_cents,stock_qty)
VALUES
  ('BR-0001','LEGO Classic Bricks', 1999, 50),
  ('DL-0002','Wooden Doll',         1299, 30),
  ('CR-0003','Remote Car',          2999, 20),
  ('BE-0004','Teddy Bear',          1599, 40)
ON CONFLICT (sku) DO NOTHING;