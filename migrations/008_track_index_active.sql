CREATE TABLE IF NOT EXISTS index_active (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  tx_id INT REFERENCES tx(id),
  account_address VARCHAR(42) NOT NULL DEFAULT '0x0000000000000000000000000000000000000000',
  active BOOLEAN NOT NULL,
  contract_address VARCHAR(42) NOT NULL DEFAULT '0x0000000000000000000000000000000000000000',
  log_index BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_index_active_unique
  ON index_active(tx_id, account_address, active, contract_address, log_index);

CREATE INDEX idx_index_active_tx_id
  ON index_active(tx_id);

CREATE INDEX idx_index_active_account
  ON index_active(account_address, contract_address);
