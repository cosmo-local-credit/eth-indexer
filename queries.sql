--name: insert-token-transfer
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: sender_address
-- $6: recipient_address
-- $7: transfer_value
-- $8: contract_address
-- $9: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO token_transfer(tx_id, sender_address, recipient_address, transfer_value, contract_address, log_index)
SELECT id, $5, $6, $7, $8, $9 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-token-mint
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: minter_address
-- $6: recipient_address
-- $7: mint_value
-- $8: contract_address
-- $9: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO token_mint(tx_id, minter_address, recipient_address, mint_value, contract_address, log_index)
SELECT id, $5, $6, $7, $8, $9 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-token-burn
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: burner_address
-- $6: burn_value
-- $7: contract_address
-- $8: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO token_burn(tx_id, burner_address, burn_value, contract_address, log_index)
SELECT id, $5, $6, $7, $8 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-faucet-give
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: token_address
-- $6: recipient_address
-- $7: give_value
-- $8: contract_address
-- $9: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO faucet_give(tx_id, token_address, recipient_address, give_value, contract_address, log_index)
SELECT id, $5, $6, $7, $8, $9 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-pool-swap
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: initiator_address
-- $6: token_in_address
-- $7: token_out_address
-- $8: in_value
-- $9: out_value
-- $10: fee
-- $11: contract_address
-- $12: log_index
-- $13: quoted_out_value
-- $14: nominal_out_value
-- $15: protocol_fee
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO pool_swap(tx_id, initiator_address, token_in_address, token_out_address, in_value, out_value, fee, contract_address, log_index, quoted_out_value, nominal_out_value, protocol_fee)
SELECT id, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-pool-deposit
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: initiator_address
-- $6: token_in_address
-- $7: in_value
-- $8: contract_address
-- $9: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO pool_deposit(tx_id, initiator_address, token_in_address, in_value, contract_address, log_index)
SELECT id, $5, $6, $7, $8, $9 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-ownership-change
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: previous_owner
-- $6: new_owner
-- $7: contract_address
-- $8: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO ownership_change(tx_id, previous_owner, new_owner, contract_address, log_index)
SELECT id, $5, $6, $7, $8 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-index-active
-- $1: tx_hash
-- $2: block_number
-- $3: date_block
-- $4: success
-- $5: account_address
-- $6: active
-- $7: contract_address
-- $8: log_index
WITH upsert_tx AS (
    INSERT INTO tx(tx_hash, block_number, date_block, success)
    VALUES($1, $2, $3, $4)
    ON CONFLICT (tx_hash) DO UPDATE SET tx_hash = EXCLUDED.tx_hash
    RETURNING id
)
INSERT INTO index_active(tx_id, account_address, active, contract_address, log_index)
SELECT id, $5, $6, $7, $8 FROM upsert_tx
ON CONFLICT DO NOTHING

--name: insert-token
-- $1: contract_address
-- $2: token_name
-- $3: token_symbol
-- $4: token_decimals
-- $5: sink_address
INSERT INTO tokens(
	contract_address,
	token_name,
	token_symbol,
	token_decimals,
    sink_address
) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING

--name: insert-pool
-- $1: contract_address
-- $2: pool_name
-- $3: pool_symbol
INSERT INTO pools(
	contract_address,
    pool_name,
    pool_symbol
) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING

--name: restore-pool
-- $1: contract_address
UPDATE pools SET removed = false WHERE contract_address = $1

--name: restore-token
-- $1: contract_address
UPDATE tokens SET removed = false WHERE contract_address = $1

--name: remove-pool
-- $1: contract_address
UPDATE pools SET removed = true WHERE contract_address = $1

--name: remove-token
-- $1: contract_address
UPDATE tokens SET removed = true WHERE contract_address = $1
