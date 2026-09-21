-- 1. Most Critical: Date filtering on tx table
CREATE INDEX idx_tx_date_block_success
ON tx(date_block, success)
WHERE success = true;

-- 2. Critical: Join optimization
CREATE INDEX idx_token_transfer_tx_id
ON token_transfer(tx_id);

-- 3. Important: Contract filtering
CREATE INDEX idx_token_transfer_contract_address
ON token_transfer(contract_address);

-- 4. Important: Composite for common pattern
CREATE INDEX idx_token_transfer_contract_tx
ON token_transfer(contract_address, tx_id);

-- 5. For COUNT(DISTINCT sender)
CREATE INDEX idx_token_transfer_sender
ON token_transfer(sender_address);

-- 6. Composite for sender queries
CREATE INDEX idx_token_transfer_sender_contract
ON token_transfer(sender_address, contract_address);

-- 7. For pool operations
CREATE INDEX idx_pool_swap_tx_id
ON pool_swap(tx_id);

CREATE INDEX idx_pool_deposit_tx_id
ON pool_deposit(tx_id);

-- 8. For token mint/burn (if heavily used)
CREATE INDEX idx_token_mint_tx_id
ON token_mint(tx_id);

CREATE INDEX idx_token_burn_tx_id
ON token_burn(tx_id);