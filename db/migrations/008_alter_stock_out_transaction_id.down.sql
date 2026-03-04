ALTER TABLE stock_out 
ALTER COLUMN transaction_id TYPE UUID USING transaction_id::UUID;
