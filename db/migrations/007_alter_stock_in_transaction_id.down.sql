ALTER TABLE stock_in 
ALTER COLUMN transaction_id TYPE UUID USING transaction_id::UUID;
