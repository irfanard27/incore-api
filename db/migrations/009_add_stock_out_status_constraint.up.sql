ALTER TABLE stock_out 
ADD CONSTRAINT chk_stock_out_status 
CHECK (status IN ('DRAFT', 'ALLOCATED', 'IN_PROGRESS', 'DONE', 'CANCELLED'));
