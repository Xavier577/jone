// Package config provides configuration types for the jone migration tool.
package config

import "fmt"

// Config holds the main configuration for jone.
type Config struct {
	Client     string
	Connection Connection
	Migrations Migrations
}

// Connection holds database connection parameters.
type Connection struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string // disable, require, verify-full
}

// DSN returns the PostgreSQL connection string.
func (c *Connection) DSN() string {
	sslMode := c.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, sslMode)
}

// Migrations holds migration-specific configuration.
type Migrations struct {
	TableName string
}
