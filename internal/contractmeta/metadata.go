package contractmeta

import (
	"context"
	"errors"
	"log/slog"

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

type TokenMetadata struct {
	ContractAddress string
	Name            string
	Symbol          string
	Decimals        uint8
	SinkAddress     string
}

type PoolMetadata struct {
	ContractAddress string
	Name            string
	Symbol          string
}

func LoadToken(ctx context.Context, provider *ethutils.Provider, address string) (TokenMetadata, error) {
	return LoadTokenWithLogger(ctx, provider, address, nil)
}

func LoadTokenWithLogger(ctx context.Context, provider *ethutils.Provider, address string, logg *slog.Logger) (TokenMetadata, error) {
	var (
		name        string
		symbol      string
		decimals    uint8
		sinkAddress common.Address

		batchErr w3.CallErrors
	)

	if logg != nil {
		logg.Debug("fetching token metadata", "address", address)
	}

	contractAddress := w3.A(address)
	if err := provider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, nameGetter).Returns(&name),
		eth.CallFunc(contractAddress, symbolGetter).Returns(&symbol),
		eth.CallFunc(contractAddress, decimalsGetter).Returns(&decimals),
	); errors.As(err, &batchErr) {
		return TokenMetadata{}, batchErr
	} else if err != nil {
		return TokenMetadata{}, err
	}

	if logg != nil {
		logg.Debug(
			"fetched token name symbol decimals",
			"address", address,
			"name", name,
			"symbol", symbol,
			"decimals", decimals,
		)
	}

	if logg != nil {
		logg.Debug("fetching token sink address", "address", address)
	}
	if err := provider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, sinkAddressGetter).Returns(&sinkAddress),
	); err != nil {
		sinkAddress = ethutils.ZeroAddress
	}

	if logg != nil {
		logg.Debug("fetched token sink address", "address", address, "sink_address", sinkAddress.Hex())
	}

	return TokenMetadata{
		ContractAddress: common.HexToAddress(address).Hex(),
		Name:            name,
		Symbol:          symbol,
		Decimals:        decimals,
		SinkAddress:     sinkAddress.Hex(),
	}, nil
}

func LoadPool(ctx context.Context, provider *ethutils.Provider, address string) (PoolMetadata, error) {
	return LoadPoolWithLogger(ctx, provider, address, nil)
}

func LoadPoolWithLogger(ctx context.Context, provider *ethutils.Provider, address string, logg *slog.Logger) (PoolMetadata, error) {
	var (
		name   string
		symbol string
	)

	if logg != nil {
		logg.Debug("fetching pool metadata", "address", address)
	}

	contractAddress := w3.A(address)
	if err := provider.Client.CallCtx(
		ctx,
		eth.CallFunc(contractAddress, nameGetter).Returns(&name),
		eth.CallFunc(contractAddress, symbolGetter).Returns(&symbol),
	); err != nil {
		return PoolMetadata{}, err
	}

	if logg != nil {
		logg.Debug(
			"fetched pool metadata",
			"address", address,
			"name", name,
			"symbol", symbol,
		)
	}

	return PoolMetadata{
		ContractAddress: common.HexToAddress(address).Hex(),
		Name:            name,
		Symbol:          symbol,
	}, nil
}
