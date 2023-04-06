package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Generic type containing specific configuration
type Config[T any] struct {
	Config *T
}

func NewConfig[T any]() *Config[T] {
	return &Config[T]{new(T)}
}

// Envelope for proceeding main configuration struct
type config[T any] struct {
	object *Config[T]
	reader io.Reader
	Error  error
}

// Create envelope from reader
func (this *Config[T]) from(reader io.Reader, err error) *config[T] {
	return &config[T]{
		object: this,
		reader: reader,
		Error:  err,
	}
}

// Set file
func (this *Config[T]) FromFile(name string) *config[T] {
	return this.from(os.Open(name))
}

// Set reader
func (this *Config[T]) FromReader(reader io.Reader) *config[T] {
	return this.from(reader, nil)
}

// Parse yaml file
func (this *config[T]) Yaml() (*T, error) {
	if this.Error != nil {
		return nil, this.Error
	}
	if this.reader == nil {
		return nil, nil
	}
	return this.object.Config, yaml.NewDecoder(this.reader).Decode(this.object.Config)
}
