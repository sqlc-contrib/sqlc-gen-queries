// Package sqlc provides functionality for parsing and working with sqlc configuration files.
//
// This package handles reading sqlc.yaml configuration files and extracting
// SQL generation settings including schema paths, query paths, and database engines.
package sqlc

import (
	"os"

	"gopkg.in/yaml.v3"
)

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
	Schema      string   `yaml:"schema"`
	Engine      string   `yaml:"engine"`
	Queries     string   `yaml:"queries"`
	SkipQueries []string `yaml:"skip_queries,omitempty"`
}

// GetSkipQueriesSet returns a set of query names to skip for efficient lookup.
// The returned map allows O(1) lookup to check if a query should be skipped.
func (s *SQL) GetSkipQueriesSet() map[string]bool {
	skipSet := make(map[string]bool, len(s.SkipQueries))
	for _, name := range s.SkipQueries {
		skipSet[name] = true
	}
	return skipSet
}
