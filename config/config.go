// Package config provides configuration types for the jone migration tool.
package config

// Config holds the main configuration for jone.
type Config struct {
	Client     string
	Connection Connection
	Migrations Migrations
}

// Connection holds database connection parameters.
type Connection struct {
	User     string
	Password string
	Database string
	Port     string
	Host     string
}

// Migrations holds migration-specific configuration.
type Migrations struct {
	TableName string
}
