package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/cosmo-local-credit/eth-indexer/internal/reconciler"
	"github.com/cosmo-local-credit/eth-indexer/internal/store"
	"github.com/cosmo-local-credit/eth-indexer/internal/util"
	"github.com/grassrootseconomics/ethutils"
	"github.com/knadh/koanf/v2"
)

const defaultReconcilerTimeout = time.Hour

var (
	build = "dev"

	confFlag             string
	migrationsFolderFlag string
	queriesFlag          string
	timeoutFlag          time.Duration

	lo *slog.Logger
	ko *koanf.Koanf
)

func init() {
	flag.StringVar(&confFlag, "config", "config.toml", "Config file location")
	flag.StringVar(&migrationsFolderFlag, "migrations", "migrations/", "Migrations folder location")
	flag.StringVar(&queriesFlag, "queries", "queries.sql", "Queries file location")
	flag.DurationVar(&timeoutFlag, "timeout", defaultReconcilerTimeout, "Reconciler timeout")
	flag.Parse()

	lo = util.InitLogger()
	ko = util.InitConfig(lo, confFlag)

	lo.Info("starting cosmo-local-credit indexer reconciler", "build", build)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutFlag)
	defer cancel()

	pgStore, err := store.NewPgStore(store.PgOpts{
		Logg:                 lo,
		DSN:                  ko.MustString("postgres.dsn"),
		MigrationsFolderPath: migrationsFolderFlag,
		QueriesFolderPath:    queriesFlag,
	})
	if err != nil {
		lo.Error("could not initialize postgres store", "error", err)
		os.Exit(1)
	}
	defer pgStore.Close()

	chainProvider := ethutils.NewProvider(
		ko.MustString("chain.rpc_endpoint"),
		ko.MustInt64("chain.chainid"),
	)

	reconcilerService := reconciler.New(reconciler.ReconcilerOpts{
		Store:         pgStore,
		ChainProvider: chainProvider,
		Logg:          lo,
	})

	result, err := reconcilerService.Run(ctx)
	if err != nil {
		lo.Error("reconciliation failed", "error", err)
		os.Exit(1)
	}

	lo.Info(
		"reconciliation complete",
		"token_scanned", result.Tokens.Scanned,
		"token_restored", result.Tokens.Restored,
		"token_insert_tried", result.Tokens.InsertTried,
		"pool_scanned", result.Pools.Scanned,
		"pool_restored", result.Pools.Restored,
		"pool_insert_tried", result.Pools.InsertTried,
	)
}
