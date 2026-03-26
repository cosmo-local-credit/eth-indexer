package handler

import (
	"log/slog"
	"strings"

	"github.com/cosmo-local-credit/eth-indexer/internal/cache"
	"github.com/cosmo-local-credit/eth-indexer/internal/store"
	"github.com/grassrootseconomics/ethutils"
)

const (
	tokenIndexContractAddress = "0xe2CEf4000d6003958c891D251328850f84654eb9"
	poolIndexContractAddress  = "0x01eD8Fe01a2Ca44Cb26D00b1309d7D777471D00C"
)

type (
	HandlerOpts struct {
		Store         store.Store
		Cache         *cache.Cache
		ChainProvider *ethutils.Provider
		Logg          *slog.Logger
	}

	Handler struct {
		store         store.Store
		cache         *cache.Cache
		chainProvider *ethutils.Provider
		logg          *slog.Logger
	}
)

func NewHandler(o HandlerOpts) *Handler {
	return &Handler{
		store:         o.Store,
		cache:         o.Cache,
		chainProvider: o.ChainProvider,
		logg:          o.Logg,
	}
}

func isFeaturedPoolIndexContractAddress(address string) bool {
	return strings.EqualFold(address, poolIndexContractAddress)
}

func isRemovableIndexContractAddress(address string) bool {
	return isFeaturedPoolIndexContractAddress(address) || strings.EqualFold(address, tokenIndexContractAddress)
}
