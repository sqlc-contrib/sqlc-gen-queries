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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SQL represents a single SQL configuration block within the sqlc.yaml file.
// Each SQL block defines the schema files, query files, and database engine
// for a specific code generation target.
type SQL struct {
	Schema  string `yaml:"schema"`
	Queries string `yaml:"queries"`
	Engine  string `yaml:"engine"`
}
