# sqlc-gen-queries

[![CI](https://github.com/sqlc-contrib/sqlc-gen-queries/actions/workflows/ci.yml/badge.svg)](https://github.com/sqlc-contrib/sqlc-gen-queries/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/sqlc-contrib/sqlc-gen-queries?include_prereleases)](https://github.com/sqlc-contrib/sqlc-gen-queries/releases)
[![License](https://img.shields.io/github/license/sqlc-contrib/sqlc-gen-queries)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![sqlc](https://img.shields.io/badge/sqlc-compatible-blue)](https://sqlc.dev)

A CLI tool that generates [sqlc](https://sqlc.dev)-compatible SQL queries from your database schema catalog. Point it at a schema catalog and a configuration file, and it produces ready-to-use query files for sqlc.

## Features

- Generate CRUD queries (`SELECT`, `INSERT`, `UPDATE`, `DELETE`) from a database schema catalog
- Configurable via YAML — control which tables and operations to generate
- Works as a standalone CLI or as part of a CI/CD pipeline
- Supports custom query templates

## Installation

```bash
go install github.com/sqlc-contrib/sqlc-gen-queries/cmd/sqlc-gen-queries@latest
```

## Usage

```bash
sqlc-gen-queries --config-file sqlc.yaml --catalog-file schema.json
```

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--config-file` | `SQLC_CONFIG_FILE` | | Path to the sqlc configuration file |
| `--catalog-file` | `SQLC_CATALOG_FILE` | `schema.json` | Path to the catalog file |

## Contributing

Contributions are welcome! Please open an issue or pull request.

To set up a development environment with [Nix](https://nixos.org):

```bash
nix develop
go test ./...
```

## License

[MIT](LICENSE)
