package service

// Configuration for service
type Config struct {
	WisdomFile string `yaml:"wisdom_file"`
	Difficulty int32  `yaml:"difficulty"`
	Port       string `yaml:"port"`
	Multicore  bool   `yaml:"multicore"`
}
