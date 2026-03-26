package handler

import (
	"context"
	"errors"

	"github.com/cosmo-local-credit/eth-tracker/pkg/event"
	"github.com/ethereum/go-ethereum/common"
	"github.com/grassrootseconomics/ethutils"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
)

var (
	nameGetter        = w3.MustNewFunc("name()", "string")
	symbolGetter      = w3.MustNewFunc("symbol()", "string")
	decimalsGetter    = w3.MustNewFunc("decimals()", "uint8")
	sinkAddressGetter = w3.MustNewFunc("sinkAddress()", "address")
)

func (h *Handler) AddToken(ctx context.Context, event event.Event) error {
	if h.cache.Get(event.ContractAddress) {
		return nil
	}

	var (
		tokenName     string
		tokenSymbol   string
		tokenDecimals uint8
		sinkAddress   common.Address

		batchErr w3.CallErrors
	)

	contractAddress := w3.A(event.ContractAddress)

	if err := h.chainProvider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, nameGetter).Returns(&tokenName),
		eth.CallFunc(contractAddress, symbolGetter).Returns(&tokenSymbol),
		eth.CallFunc(contractAddress, decimalsGetter).Returns(&tokenDecimals),
	); errors.As(err, &batchErr) {
		return batchErr
	} else if err != nil {
		return err
	}

	if err := h.chainProvider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, sinkAddressGetter).Returns(&sinkAddress),
	); err != nil {
		// This will most likely revert if the contract does not have a sinkAddress
		// Instead of handling the error we just ignore it and set the value to 0
		sinkAddress = ethutils.ZeroAddress
	}

	if err := h.store.InsertToken(ctx, event.ContractAddress, tokenName, tokenSymbol, tokenDecimals, sinkAddress.Hex()); err != nil {
		return err
	}

	h.cache.Set(event.ContractAddress)
	return nil
}

func (h *Handler) AddPool(ctx context.Context, event event.Event) error {
	if h.cache.Get(event.ContractAddress) {
		return nil
	}

	var (
		tokenName   string
		tokenSymbol string
	)

	contractAddress := w3.A(event.ContractAddress)

	if err := h.chainProvider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, nameGetter).Returns(&tokenName),
		eth.CallFunc(contractAddress, symbolGetter).Returns(&tokenSymbol),
	); err != nil {
		return err
	}

	if err := h.store.InsertPool(ctx, event.ContractAddress, tokenName, tokenSymbol); err != nil {
		return err
	}

	h.cache.Set(event.ContractAddress)
	return nil
}

// This is a special method meant to improve the UX on https://sarafu.network/pools
func (h *Handler) AddSarafuNetworkFeaturedPool(ctx context.Context, event event.Event) error {
	// This is the only pool index
	if !isFeaturedPoolIndexContractAddress(event.ContractAddress) {
		return nil
	}

	var (
		tokenName   string
		tokenSymbol string
	)

	poolAddress := event.Payload["address"].(string)
	contractAddress := w3.A(poolAddress)

	if err := h.chainProvider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, nameGetter).Returns(&tokenName),
		eth.CallFunc(contractAddress, symbolGetter).Returns(&tokenSymbol),
	); err != nil {
		return err
	}

	if err := h.store.InsertPool(ctx, poolAddress, tokenName, tokenSymbol); err != nil {
		return err
	}

	h.cache.Set(poolAddress)
	return nil
}
