CREATE TABLE stock_out_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_out_id UUID NOT NULL,
    inventory_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (stock_out_id) REFERENCES stock_out(id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE CASCADE
);

-- create indexes
CREATE INDEX idx_stock_out_items_stock_out_id ON stock_out_items(stock_out_id);
CREATE INDEX idx_stock_out_items_inventory_id ON stock_out_items(inventory_id);
