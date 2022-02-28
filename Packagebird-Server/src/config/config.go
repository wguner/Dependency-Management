package config

// Global configuration structure
var Config *ConfigFile

// Structure type
type ConfigFile struct {
	PackageSourcePath string
	ProjectSourcePath string
	ServerAddress     string
	DatabaseAddress   string
}
