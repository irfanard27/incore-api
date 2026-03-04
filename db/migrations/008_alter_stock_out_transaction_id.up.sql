ALTER TABLE stock_out 
ALTER COLUMN transaction_id TYPE VARCHAR(255) USING transaction_id::VARCHAR(255);
