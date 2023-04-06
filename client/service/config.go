package service

// Configuration for client
type Config struct {
	Address   string `yaml:"address"`
	Multicore bool   `yaml:"multicore"`
}
