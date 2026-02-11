
-- WARNING: This specific migration will also delete transfers from the same recipient, same sender, same value in the same tx and reduce it to 1 because we are not collecting event logs to distinguish them
DELETE FROM token_transfer
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, sender_address, recipient_address, transfer_value, contract_address
            ORDER BY id
        ) AS rn
        FROM token_transfer
    ) t WHERE t.rn > 1
);

DELETE FROM token_mint
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, minter_address, recipient_address, mint_value, contract_address
            ORDER BY id
        ) AS rn
        FROM token_mint
    ) t WHERE t.rn > 1
);

DELETE FROM token_burn
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, burner_address, burn_value, contract_address
            ORDER BY id
        ) AS rn
        FROM token_burn
    ) t WHERE t.rn > 1
);

DELETE FROM faucet_give
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, token_address, recipient_address, give_value, contract_address
            ORDER BY id
        ) AS rn
        FROM faucet_give
    ) t WHERE t.rn > 1
);

DELETE FROM pool_swap
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, initiator_address, token_in_address, token_out_address, in_value, out_value, fee, contract_address
            ORDER BY id
        ) AS rn
        FROM pool_swap
    ) t WHERE t.rn > 1
);

DELETE FROM pool_deposit
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, initiator_address, token_in_address, in_value, contract_address
            ORDER BY id
        ) AS rn
        FROM pool_deposit
    ) t WHERE t.rn > 1
);

DELETE FROM ownership_change
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (
            PARTITION BY tx_id, previous_owner, new_owner, contract_address
            ORDER BY id
        ) AS rn
        FROM ownership_change
    ) t WHERE t.rn > 1
);


CREATE UNIQUE INDEX idx_token_transfer_unique
ON public.token_transfer(tx_id, sender_address, recipient_address, transfer_value, contract_address);

CREATE UNIQUE INDEX idx_token_mint_unique
ON public.token_mint(tx_id, minter_address, recipient_address, mint_value, contract_address);

CREATE UNIQUE INDEX idx_token_burn_unique
ON public.token_burn(tx_id, burner_address, burn_value, contract_address);

CREATE UNIQUE INDEX idx_faucet_give_unique
ON public.faucet_give(tx_id, token_address, recipient_address, give_value, contract_address);

CREATE UNIQUE INDEX idx_pool_swap_unique
ON public.pool_swap(tx_id, initiator_address, token_in_address, token_out_address, in_value, out_value, fee, contract_address);

CREATE UNIQUE INDEX idx_pool_deposit_unique
ON public.pool_deposit(tx_id, initiator_address, token_in_address, in_value, contract_address);

CREATE UNIQUE INDEX idx_ownership_change_unique
ON public.ownership_change(tx_id, previous_owner, new_owner, contract_address);
