package handler

import (
	"context"

	"github.com/cosmo-local-credit/eth-indexer/internal/indexcontracts"
	"github.com/cosmo-local-credit/eth-tracker/pkg/event"
)

func (h *Handler) IndexAdd(ctx context.Context, event event.Event) error {
	if indexcontracts.IsTrackedIndexContractAddress(event.ContractAddress) {
		return h.store.RestoreContractAddress(ctx, event)
	}

	return nil
}
