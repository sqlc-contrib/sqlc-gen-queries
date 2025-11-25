package sqlc

// Generator represents the SQLC Queries generator.
type Generator struct {
	Config  *Config
	Catalog *Catalog
}

// Generate generates the queries based on the configuration.
func (x *Generator) Generate() error {
	return nil
}
