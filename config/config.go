package config

type Config struct {
	Port             *int
	Shipnode         *string
	ConnectionString string
}

func New(port *int, shipnode *string, db string) *Config {
	return &Config{
		port,
		shipnode,
		db,
	}
}
