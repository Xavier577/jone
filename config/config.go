// Package config provides configuration types for the jone migration tool.
package config

import (
	"fmt"
	"time"
)

// Config holds the main configuration for jone.
type Config struct {
	Client     string
	Connection Connection
	Pool       Pool
	Migrations Migrations
}

// Pool holds connection pool configuration.
// Zero values preserve database/sql defaults.
type Pool struct {
	MaxOpenConns    int           // Maximum number of open connections. 0 means unlimited.
	MaxIdleConns    int           // Maximum number of idle connections. 0 means default (2).
	ConnMaxLifetime time.Duration // Maximum time a connection may be reused. 0 means no limit.
	ConnMaxIdleTime time.Duration // Maximum time a connection may be idle. 0 means no limit.
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
//
// Deprecated: Use Dialect.FormatDSN instead, which handles all database types.
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
