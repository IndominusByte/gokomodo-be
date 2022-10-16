DROP TABLE IF EXISTS transaction.orders;
DROP INDEX IF EXISTS idx_transaction_orders_buyer_id ON transaction.orders(buyer_id);
DROP INDEX IF EXISTS idx_transaction_orders_seller_id ON transaction.orders(seller_id);
DROP INDEX IF EXISTS idx_transaction_orders_product_id ON transaction.orders(product_id);
DROP INDEX IF EXISTS idx_transaction_orders_status ON transaction.orders(status);
