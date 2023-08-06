package database

import "context"

type Database interface {
	// Close closes the connection to storage.
	Close() error
	// Ping - checks if storage is available.
	Ping(ctx context.Context) error
}
