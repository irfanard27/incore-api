CREATE TABLE stock_in_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_in_id UUID NOT NULL,
    inventory_id UUID NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (stock_in_id) REFERENCES stock_in(id) ON DELETE CASCADE,
    FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE CASCADE
);

-- create indexes
CREATE INDEX idx_stock_in_items_stock_in_id ON stock_in_items(stock_in_id);
CREATE INDEX idx_stock_in_items_inventory_id ON stock_in_items(inventory_id);
