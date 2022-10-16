CREATE TABLE IF NOT EXISTS transaction.orders(
  id SERIAL PRIMARY KEY,
  buyer_id INT NOT NULL,
  seller_id INT NOT NULL,
  product_id INT NOT NULL,
  source_address TEXT NOT NULL,
  destination_address TEXT NOT NULL,
  qty BIGINT NOT NULL,
  price BIGINT NOT NULL,
  total_price BIGINT NOT NULL,
  status VARCHAR(20) DEFAULT 'pending',
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transaction_orders_buyer_id ON transaction.orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_transaction_orders_seller_id ON transaction.orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_transaction_orders_product_id ON transaction.orders(product_id);
CREATE INDEX IF NOT EXISTS idx_transaction_orders_status ON transaction.orders(status);
