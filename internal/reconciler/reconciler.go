package reconciler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cosmo-local-credit/eth-indexer/internal/contractmeta"
	"github.com/cosmo-local-credit/eth-indexer/internal/indexcontracts"
	"github.com/cosmo-local-credit/eth-indexer/internal/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/grassrootseconomics/ethutils"
)

type ReconcilerOpts struct {
	Store         store.Store
	ChainProvider *ethutils.Provider
	Logg          *slog.Logger
}

type Reconciler struct {
	store         store.Store
	chainProvider *ethutils.Provider
	logg          *slog.Logger
}

type Result struct {
	Tokens ScannedResult
	Pools  ScannedResult
}

type ScannedResult struct {
	Scanned     int
	Restored    int64
	InsertTried int
	SkippedZero int
}

func New(o ReconcilerOpts) *Reconciler {
	return &Reconciler{
		store:         o.Store,
		chainProvider: o.ChainProvider,
		logg:          o.Logg,
	}
}

func (r *Reconciler) Run(ctx context.Context) (Result, error) {
	result := Result{}
	r.logg.Info("starting token reconciliation", "index_contract", indexcontracts.TokenContractAddress)

	tokens, err := r.reconcileTokens(ctx)
	if err != nil {
		return Result{}, err
	}
	result.Tokens = tokens
	r.logg.Info(
		"token reconciliation finished",
		"scanned", tokens.Scanned,
		"restored", tokens.Restored,
		"insert_tried", tokens.InsertTried,
		"skipped_zero", tokens.SkippedZero,
	)

	r.logg.Info("starting pool reconciliation", "index_contract", indexcontracts.PoolContractAddress)

	pools, err := r.reconcilePools(ctx)
	if err != nil {
		return Result{}, err
	}
	result.Pools = pools
	r.logg.Info(
		"pool reconciliation finished",
		"scanned", pools.Scanned,
		"restored", pools.Restored,
		"insert_tried", pools.InsertTried,
		"skipped_zero", pools.SkippedZero,
	)

	return result, nil
}

func (r *Reconciler) reconcileTokens(ctx context.Context) (ScannedResult, error) {
	r.logg.Debug("creating token index iterator", "index_contract", indexcontracts.TokenContractAddress)
	iter, err := r.chainProvider.NewBatchIterator(ctx, common.HexToAddress(indexcontracts.TokenContractAddress))
	if err != nil {
		return ScannedResult{}, fmt.Errorf("create token index iterator: %w", err)
	}
	r.logg.Debug("created token index iterator", "index_contract", indexcontracts.TokenContractAddress)

	return r.reconcileIndexBatches(
		ctx,
		"token",
		iter,
		r.store.RestoreTokensByAddress,
		func(ctx context.Context, address string) error {
			metadata, err := contractmeta.LoadTokenWithLogger(ctx, r.chainProvider, address, r.logg)
			if err != nil {
				return err
			}

			r.logg.Debug("inserting token metadata", "address", metadata.ContractAddress)
			return r.store.InsertToken(
				ctx,
				metadata.ContractAddress,
				metadata.Name,
				metadata.Symbol,
				metadata.Decimals,
				metadata.SinkAddress,
			)
		},
	)
}

func (r *Reconciler) reconcilePools(ctx context.Context) (ScannedResult, error) {
	r.logg.Debug("creating pool index iterator", "index_contract", indexcontracts.PoolContractAddress)
	iter, err := r.chainProvider.NewBatchIterator(ctx, common.HexToAddress(indexcontracts.PoolContractAddress))
	if err != nil {
		return ScannedResult{}, fmt.Errorf("create pool index iterator: %w", err)
	}
	r.logg.Debug("created pool index iterator", "index_contract", indexcontracts.PoolContractAddress)

	return r.reconcileIndexBatches(
		ctx,
		"pool",
		iter,
		r.store.RestorePoolsByAddress,
		func(ctx context.Context, address string) error {
			metadata, err := contractmeta.LoadPoolWithLogger(ctx, r.chainProvider, address, r.logg)
			if err != nil {
				return err
			}

			r.logg.Debug("inserting pool metadata", "address", metadata.ContractAddress)
			return r.store.InsertPool(
				ctx,
				metadata.ContractAddress,
				metadata.Name,
				metadata.Symbol,
			)
		},
	)
}

func (r *Reconciler) reconcileIndexBatches(
	ctx context.Context,
	kind string,
	iter *ethutils.BatchIterator,
	restore func(context.Context, []string) (int64, error),
	insert func(context.Context, string) error,
) (ScannedResult, error) {
	result := ScannedResult{}
	batchNumber := 0

	for {
		batchNumber++
		r.logg.Debug("requesting index batch", "kind", kind, "batch_number", batchNumber)
		batch, err := iter.Next(ctx)
		if err != nil {
			return ScannedResult{}, fmt.Errorf("iterate %s index: %w", kind, err)
		}
		if batch == nil {
			r.logg.Info("index iteration complete", "kind", kind, "batches", batchNumber-1)
			break
		}
		r.logg.Debug("received raw index batch", "kind", kind, "batch_number", batchNumber, "raw_batch_size", len(batch))

		addresses, skippedZero := normalizeBatch(batch)
		result.SkippedZero += skippedZero
		result.Scanned += len(addresses)
		r.logg.Info(
			"processing reconciler batch",
			"kind", kind,
			"batch_number", batchNumber,
			"addresses", len(addresses),
			"skipped_zero", skippedZero,
		)
		if len(addresses) == 0 {
			continue
		}

		insertedBatch := 0
		for index, address := range addresses {
			r.logg.Debug(
				"processing contract",
				"kind", kind,
				"batch_number", batchNumber,
				"position", index+1,
				"batch_size", len(addresses),
				"address", address,
			)
			if err := insert(ctx, address); err != nil {
				return ScannedResult{}, fmt.Errorf("insert %s %s: %w", kind, address, err)
			}
			result.InsertTried++
			insertedBatch++
		}

		r.logg.Debug("restoring batch addresses", "kind", kind, "batch_number", batchNumber, "addresses", len(addresses))
		restored, err := restore(ctx, addresses)
		if err != nil {
			return ScannedResult{}, fmt.Errorf("restore %s addresses: %w", kind, err)
		}
		result.Restored += restored

		r.logg.Info(
			"reconciled index batch",
			"kind", kind,
			"batch_size", len(addresses),
			"insert_tried", insertedBatch,
			"restored", restored,
		)
	}

	return result, nil
}

func normalizeBatch(batch []common.Address) ([]string, int) {
	addresses := make([]string, 0, len(batch))
	seen := make(map[string]struct{}, len(batch))
	skippedZero := 0

	for _, address := range batch {
		if address == ethutils.ZeroAddress {
			skippedZero++
			continue
		}

		normalized := address.Hex()
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		addresses = append(addresses, normalized)
	}

	return addresses, skippedZero
}
