package main

import (
	"github.com/cosmo-local-credit/eth-indexer/internal/handler"
	"github.com/cosmo-local-credit/eth-indexer/pkg/router"
)

func bootstrapRouter(handlerContainer *handler.Handler) *router.Router {
	router := router.New(lo)

	router.RegisterRoute(
		"TRACKER.TOKEN_TRANSFER",
		handlerContainer.IndexTransfer,
		handlerContainer.AddToken,
	)
	router.RegisterRoute(
		"TRACKER.TOKEN_TRANSFER_FROM",
		handlerContainer.IndexTransfer,
		handlerContainer.AddToken,
	)
	router.RegisterRoute(
		"TRACKER.TOKEN_MINT",
		handlerContainer.IndexTokenMint,
		handlerContainer.AddToken,
	)
	router.RegisterRoute(
		"TRACKER.TOKEN_BURN",
		handlerContainer.IndexTokenBurn,
		handlerContainer.AddToken,
	)
	router.RegisterRoute(
		"TRACKER.POOL_SWAP",
		handlerContainer.IndexPoolSwap,
		handlerContainer.AddPool,
	)
	router.RegisterRoute(
		"TRACKER.POOL_DEPOSIT",
		handlerContainer.IndexPoolDeposit,
		handlerContainer.AddPool,
	)
	router.RegisterRoute(
		"TRACKER.FAUCET_GIVE",
		handlerContainer.IndexFaucetGive,
		handlerContainer.FaucetHealthCheck,
	)
	router.RegisterRoute(
		"TRACKER.OWNERSHIP_TRANSFERRED",
		handlerContainer.IndexOwnershipChange,
	)

	router.RegisterRoute(
		"TRACKER.INDEX_REMOVE",
		handlerContainer.IndexRemove,
	)

	// This is a special method meant to improve the UX on https://sarafu.network/pools
	if ko.Bool("sarafu_network.featured_pools_enabled") {
		lo.Info("sarafu network featured pools enabled, registering extra pool add handler")
		router.RegisterRoute(
			"TRACKER.INDEX_ADD",
			handlerContainer.AddSarafuNetworkFeaturedPool,
		)
	}

	return router
}
