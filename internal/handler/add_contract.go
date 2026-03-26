package handler

import (
	"context"

	"github.com/cosmo-local-credit/eth-indexer/internal/contractmeta"
	"github.com/cosmo-local-credit/eth-indexer/internal/indexcontracts"
	"github.com/cosmo-local-credit/eth-tracker/pkg/event"
)

func (h *Handler) AddToken(ctx context.Context, event event.Event) error {
	if h.cache.Get(event.ContractAddress) {
		return nil
	}

	metadata, err := contractmeta.LoadToken(ctx, h.chainProvider, event.ContractAddress)
	if err != nil {
		return err
	}

	if err := h.store.InsertToken(ctx, metadata.ContractAddress, metadata.Name, metadata.Symbol, metadata.Decimals, metadata.SinkAddress); err != nil {
		return err
	}

	h.cache.Set(metadata.ContractAddress)
	return nil
}

func (h *Handler) AddPool(ctx context.Context, event event.Event) error {
	if h.cache.Get(event.ContractAddress) {
		return nil
	}

	metadata, err := contractmeta.LoadPool(ctx, h.chainProvider, event.ContractAddress)
	if err != nil {
		return err
	}

	if err := h.store.InsertPool(ctx, metadata.ContractAddress, metadata.Name, metadata.Symbol); err != nil {
		return err
	}

	h.cache.Set(metadata.ContractAddress)
	return nil
}

// This is a special method meant to improve the UX on https://sarafu.network/pools
func (h *Handler) AddSarafuNetworkFeaturedPool(ctx context.Context, event event.Event) error {
	// This is the only pool index
	if !indexcontracts.IsPoolContractAddress(event.ContractAddress) {
		return nil
	}

	poolAddress := event.Payload["address"].(string)
	metadata, err := contractmeta.LoadPool(ctx, h.chainProvider, poolAddress)
	if err != nil {
		return err
	}

	if err := h.store.InsertPool(ctx, metadata.ContractAddress, metadata.Name, metadata.Symbol); err != nil {
		return err
	}

	h.cache.Set(metadata.ContractAddress)
	return nil
}
