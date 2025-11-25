package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/sqlc-contrib/sqlc-gen-queries/internal/sqlc"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:      "sqlc-gen-queries",
		Usage:     "SQLC Queries Generator",
		UsageText: "sqlc-gen-queries [global options] [command]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-file",
				Usage:   "Path to the sqlc configuration file.",
				Sources: cli.EnvVars("SQLC_CONFIG_FILE"),
				Value:   "sqlc.yaml",
			},
			&cli.StringFlag{
				Name:    "catalog-file",
				Usage:   "Path to the catalog file.",
				Sources: cli.EnvVars("SQLC_CATALOG_FILE"),
				Value:   "schema.json",
			},
		},
		ErrWriter: os.Stderr,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config, err := sqlc.LoadConfig(cmd.String("config-file"))
			if err != nil {
				return err
			}

			catalog, err := sqlc.LoadCatalog(cmd.String("catalog-file"))
			if err != nil {
				return err
			}

			generator := &sqlc.Generator{
				Config:  config,
				Catalog: catalog,
			}

			return generator.Generate()
		},
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		app.Version = info.Main.Version
	}

	ctx := context.Background()
	// start the application
	if err := app.Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "The generator has encountered a fatal error", slog.Any("error", err))
		os.Exit(1)
	}
}
