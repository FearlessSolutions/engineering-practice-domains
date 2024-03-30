package config

import "os"

// registrySource is a source of information for retrieving configuration options
type registrySource interface {
	getValue(variableName string) (string, bool)
}

// environmentRegistrySource is a registrySource that pulls configuration options from environment variables
type environmentRegistrySource struct{}

func (src environmentRegistrySource) getValue(variableName string) (string, bool) {
	return os.LookupEnv(variableName)
}

// mockRegistrySource is an in-memory map-based registrySource implementation appropriate for testing
type mockRegistrySource map[string]string

func (src mockRegistrySource) getValue(variableName string) (string, bool) {
	value, isPresent := src[variableName]
	return value, isPresent
}
