# Data Guide for Researchers

This database (`chain_data`) contains indexed Ethereum blockchain events.


## 1. Schema Overview

Every event row links back to a `tx` row via `tx_id`. Contract addresses in event tables can be joined to `tokens` or `pools` for human-readable metadata.


### Table Reference

#### `tx` — Transaction records

| Column | Type | Description |
|---|---|---|
| `id` | INT | Internal surrogate key. |
| `tx_hash` | VARCHAR(66) | Full Ethereum transaction hash (`0x` + 64 hex chars). Unique. |
| `block_number` | INT | Block number this transaction was included in. |
| `date_block` | TIMESTAMP | Block timestamp in **UTC**. |
| `success` | BOOLEAN | `true` if the transaction was not reverted on-chain. |

---

#### `token_transfer` — ERC-20 Transfer events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `sender_address` | VARCHAR(42) | Address that sent tokens. |
| `recipient_address` | VARCHAR(42) | Address that received tokens. |
| `contract_address` | VARCHAR(42) | Token contract. Join to `tokens` for name/symbol. |
| `transfer_value` | NUMERIC | Raw on-chain integer amount (divide by `10^token_decimals`). |
| `log_index` | BIGINT | Position of this log entry within the transaction. Used together with `tx_id` to uniquely identify the event. |

---

#### `token_mint` — Token minting events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `minter_address` | VARCHAR(42) | Address that triggered the mint (usually a privileged operator). |
| `recipient_address` | VARCHAR(42) | Address that received the newly minted tokens. |
| `contract_address` | VARCHAR(42) | Token contract. |
| `mint_value` | NUMERIC | Raw on-chain integer amount. |
| `log_index` | BIGINT | Log position within the transaction. |

---

#### `token_burn` — Token burn / sink events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `burner_address` | VARCHAR(42) | Address that burned tokens. |
| `contract_address` | VARCHAR(42) | Token contract. |
| `burn_value` | NUMERIC | Raw on-chain integer amount. |
| `log_index` | BIGINT | Log position within the transaction. |

---

#### `faucet_give` — Faucet distribution events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `token_address` | VARCHAR(42) | Address of the token being distributed. |
| `recipient_address` | VARCHAR(42) | Address that received the faucet payout. |
| `contract_address` | VARCHAR(42) | Faucet contract that emitted the event. |
| `give_value` | NUMERIC | Raw on-chain integer amount. |
| `log_index` | BIGINT | Log position within the transaction. |

---

#### `pool_swap` — Liquidity pool swap events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `initiator_address` | VARCHAR(42) | Address that initiated the swap. |
| `token_in_address` | VARCHAR(42) | Token being sold into the pool. |
| `token_out_address` | VARCHAR(42) | Token being bought out of the pool. |
| `in_value` | NUMERIC | Raw amount of `token_in` sold. |
| `out_value` | NUMERIC | Raw amount of `token_out` received. |
| `fee` | NUMERIC | Raw pool fee, excluding the protocol fee. |
| `contract_address` | VARCHAR(42) | Pool contract. Join to `pools` for name/symbol. |
| `log_index` | BIGINT | Log position within the transaction. |
| `quoted_out_value` | NUMERIC | Raw quote before fees. NULL before the settlement cutover. |
| `nominal_out_value` | NUMERIC | Raw amount the pool transferred out. NULL before the settlement cutover. |
| `protocol_fee` | NUMERIC | Raw protocol fee. NULL before the settlement cutover. |

Every value except `in_value` is denominated in `token_out`, so scale `fee`, `protocol_fee` and each out value by `token_out`'s decimals.

The last three columns are decoded from the pool's `SwapSettlement` event, which older pool implementations did not emit. They are NULL for every row indexed below the tracker's `chain.pool_settlement_block`, so treat NULL as unknown rather than zero and exclude it from fee aggregates. `nominal_out_value` and `out_value` differ only for fee-on-transfer output tokens, where `out_value` is the increase actually observed at the recipient.

---

#### `pool_deposit` — Liquidity pool deposit events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `initiator_address` | VARCHAR(42) | Address depositing into the pool. |
| `token_in_address` | VARCHAR(42) | Token being deposited. |
| `contract_address` | VARCHAR(42) | Pool contract. |
| `in_value` | NUMERIC | Raw on-chain integer amount deposited. |
| `log_index` | BIGINT | Log position within the transaction. |

---

#### `index_active` — Registry activation state changes

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `account_address` | VARCHAR(42) | Address whose activation state changed. |
| `active` | BOOLEAN | Whether the address is active as of this event. |
| `contract_address` | VARCHAR(42) | Index contract that holds the entry. |
| `log_index` | BIGINT | Log position within the transaction. |

Deactivation revokes an address from the index's `have()` authorization check while retaining the entry for enumeration and later reactivation, so a deactivated address stays out of `index_remove` and keeps its original addition time. The current state of an address is the most recent row by `tx_id` and `log_index`; an address with no rows has never been deactivated.

---

#### `ownership_change` — Contract ownership transfer events

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `tx_id` | INT | FK → `tx.id`. |
| `previous_owner` | VARCHAR(42) | Address that previously owned the contract. |
| `new_owner` | VARCHAR(42) | Address that now owns the contract. |
| `contract_address` | VARCHAR(42) | Contract whose ownership changed. |
| `log_index` | BIGINT | Log position within the transaction. |

---

#### `tokens` — Token metadata

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `contract_address` | VARCHAR(42) | Unique token contract address. |
| `token_name` | TEXT | Human-readable name (e.g. `Grassroots Economics`). |
| `token_symbol` | TEXT | Ticker symbol (e.g. `GE`). |
| `token_decimals` | INT | Decimal places. Divide raw values by `10^token_decimals` for human units. |
| `sink_address` | VARCHAR(42) | Designated sink/burn address for this token, if any. |
| `removed` | BOOLEAN | `true` if the token has been de-listed from the indexer. |

---

#### `pools` — Pool metadata

| Column | Type | Description |
|---|---|---|
| `id` | INT | Surrogate key. |
| `contract_address` | VARCHAR(42) | Unique pool contract address. |
| `pool_name` | TEXT | Human-readable pool name. |
| `pool_symbol` | TEXT | Pool ticker symbol. |
| `removed` | BOOLEAN | `true` if the pool has been de-listed from the indexer. |

---

### Conventions

- **Addresses** are checksummed Ethereu maddresses i.e. hexadecimal with `0x` prefix, 42 characters total.
- **Values** are raw on-chain integers. To convert to a human-readable amount:
  ```sql
  transfer_value / power(10, t.token_decimals)
  ```
- **Timestamps** (`date_block`) are in UTC.
- **Idempotency**: each event is stored at most once, identified by `(tx_id, ..., log_index)`. Duplicate on-chain events within the same transaction are deduplicated using the log index.
- **Removed flag**: rows in `tokens` and `pools` with `removed = true` are contracts that the indexer has stopped tracking. They are kept for historical completeness.

---

## 2. Restoring the Dump

### Option A: Docker (recommended for local analysis)

```bash
docker run -d \
  --name chain-data \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=chain_data \
  -p 5432:5432 \
  postgres:18-alpine

pg_restore -h localhost -U postgres -d chain_data /path/to/dump.pgdump

psql -h localhost -U postgres -d chain_data -f /path/to/dump.sql
```

### Option B: Existing Postgres instance

```bash
# Create the database first
createdb -h localhost -U postgres chain_data

# Restore custom-format dump
pg_restore -h localhost -U postgres -d chain_data /path/to/dump.pgdump

# Or plain SQL
psql -h localhost -U postgres -d chain_data -f /path/to/dump.sql
```

---

## 3. Querying with DuckDB

[DuckDB](https://duckdb.org/) is a fast, in-process analytical SQL engine that runs locally with no server setup. It is well suited for ad-hoc analysis of this dataset.


