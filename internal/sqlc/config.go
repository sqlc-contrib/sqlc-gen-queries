// Package sqlc provides functionality for parsing and working with sqlc configuration files.
//
// This package handles reading sqlc.yaml configuration files and extracting
// SQL generation settings including schema paths, query paths, and database engines.
package sqlc

import (
	"os"

	"gopkg.in/yaml.v3"
)

// PluginName is the name used to identify the gen-queries plugin in sqlc codegen configuration.
const PluginName = "gen-queries"

// Config represents the root structure of a sqlc.yaml configuration file.
// This file controls how sqlc generates Go code from SQL schemas and queries.
type Config struct {
	Version string `yaml:"version"`
	SQL     []SQL  `yaml:"sql"`
}

// LoadConfig loads a sqlc.yaml configuration file from the specified path
// and returns a Config struct. The file is expected to be in YAML format.
func LoadConfig(path string) (*Config, error) {
	var files []string

	if path != "" {
		files = append(files, path)
	} else {
		files = []string{"sqlc.yaml", "sqlc.yml"}
	}

	for _, path := range files {
		data, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
		}

		if err != nil {
			return nil, err
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}

		return &config, nil
	}

	return nil, os.ErrNotExist
}

// SQL represents a single SQL configuration block within the sqlc.yaml file.
// Each SQL block defines the schema files, query files, and database engine
// for a specific code generation target.
type SQL struct {
	Schema  string    `yaml:"schema"`
	Engine  string    `yaml:"engine"`
	Queries string    `yaml:"queries"`
	Codegen []Codegen `yaml:"codegen,omitempty"`
}

// Codegen represents a code generation plugin configuration block.
type Codegen struct {
	Plugin  string         `yaml:"plugin"`
	Out     string         `yaml:"out"`
	Options CodegenOptions `yaml:"options,omitempty"`
}

// CodegenOptions holds plugin-specific options for the gen-queries plugin.
type CodegenOptions struct {
	Queries []string     `yaml:"queries,omitempty"`
	Tables  TableOptions `yaml:"tables,omitempty"`
}

// TableOptions holds table-level filtering options for the gen-queries plugin.
// Include is an allow-list: when non-empty, only the listed tables are
// generated. Exclude is a deny-list that always takes precedence over Include.
type TableOptions struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

// GetOptions returns the CodegenOptions for the gen-queries plugin.
// If no matching codegen entry is found, returns an empty CodegenOptions.
func (s *SQL) GetOptions() CodegenOptions {
	for _, c := range s.Codegen {
		if c.Plugin == PluginName {
			return c.Options
		}
	}
	return CodegenOptions{}
}

// GetQueriesSet returns a set of opt-in query names for efficient lookup.
// The returned map allows O(1) lookup to check if a query is enabled.
func (s *SQL) GetQueriesSet() map[string]bool {
	opts := s.GetOptions()
	querySet := make(map[string]bool, len(opts.Queries))
	for _, name := range opts.Queries {
		querySet[name] = true
	}
	return querySet
}

// GetIncludeSet returns the allow-list of table names for query generation.
// Entries may be unqualified table names or schema-qualified names. An empty
// set means every table is included.
func (s *SQL) GetIncludeSet() map[string]bool {
	opts := s.GetOptions()
	includeSet := make(map[string]bool, len(opts.Tables.Include))
	for _, name := range opts.Tables.Include {
		includeSet[name] = true
	}
	return includeSet
}

// GetExcludeSet returns the deny-list of table names to skip during query
// generation. Entries may be unqualified table names or schema-qualified names.
func (s *SQL) GetExcludeSet() map[string]bool {
	opts := s.GetOptions()
	excludeSet := make(map[string]bool, len(opts.Tables.Exclude))
	for _, name := range opts.Tables.Exclude {
		excludeSet[name] = true
	}
	return excludeSet
}

// tableSelected reports whether a table should have query files generated.
// Exclude always takes precedence over include; an empty include set matches
// every table. Both sets are checked against the unqualified table name and
// the schema-qualified name (schema.table).
func tableSelected(includeSet, excludeSet map[string]bool, schema, table string) bool {
	qualified := schema + "." + table
	if excludeSet[table] || excludeSet[qualified] {
		return false
	}
	if len(includeSet) == 0 {
		return true
	}
	return includeSet[table] || includeSet[qualified]
}
