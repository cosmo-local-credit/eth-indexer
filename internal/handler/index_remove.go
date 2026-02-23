package handler

import (
	"context"

	"github.com/cosmo-local-credit/eth-tracker/pkg/event"
)

func (h *Handler) IndexRemove(ctx context.Context, event event.Event) error {
	return h.store.RemoveContractAddress(ctx, event)
}
