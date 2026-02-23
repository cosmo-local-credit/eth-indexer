-- Drop old unique indexes (created in migration 005)
DROP INDEX IF EXISTS idx_token_transfer_unique;
DROP INDEX IF EXISTS idx_token_mint_unique;
DROP INDEX IF EXISTS idx_token_burn_unique;
DROP INDEX IF EXISTS idx_faucet_give_unique;
DROP INDEX IF EXISTS idx_pool_swap_unique;
DROP INDEX IF EXISTS idx_pool_deposit_unique;
DROP INDEX IF EXISTS idx_ownership_change_unique;

-- Add log_index column to all event tables (default 0 for existing rows)
ALTER TABLE token_transfer   ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;
ALTER TABLE token_mint       ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;
ALTER TABLE token_burn       ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;
ALTER TABLE faucet_give      ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;
ALTER TABLE pool_swap        ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;
ALTER TABLE pool_deposit     ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;
ALTER TABLE ownership_change ADD COLUMN IF NOT EXISTS log_index BIGINT NOT NULL DEFAULT 0;

-- Re-create unique indexes with log_index included
CREATE UNIQUE INDEX idx_token_transfer_unique
  ON token_transfer(tx_id, sender_address, recipient_address, transfer_value, contract_address, log_index);

CREATE UNIQUE INDEX idx_token_mint_unique
  ON token_mint(tx_id, minter_address, recipient_address, mint_value, contract_address, log_index);

CREATE UNIQUE INDEX idx_token_burn_unique
  ON token_burn(tx_id, burner_address, burn_value, contract_address, log_index);

CREATE UNIQUE INDEX idx_faucet_give_unique
  ON faucet_give(tx_id, token_address, recipient_address, give_value, contract_address, log_index);

CREATE UNIQUE INDEX idx_pool_swap_unique
  ON pool_swap(tx_id, initiator_address, token_in_address, token_out_address, in_value, out_value, fee, contract_address, log_index);

CREATE UNIQUE INDEX idx_pool_deposit_unique
  ON pool_deposit(tx_id, initiator_address, token_in_address, in_value, contract_address, log_index);

CREATE UNIQUE INDEX idx_ownership_change_unique
  ON ownership_change(tx_id, previous_owner, new_owner, contract_address, log_index);
