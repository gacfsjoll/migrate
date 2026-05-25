// Package migrate provides database migration functionality.
// It is a fork of golang-migrate/migrate with additional features and fixes.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// DefaultPrefetchMigrations is the default number of migrations to prefetch.
// Increased from 10 to 20 for better performance on larger migration sets.
const DefaultPrefetchMigrations = 20

// ErrNoChange is returned when no migration is needed.
var ErrNoChange = errors.New("no change")

// ErrNilVersion is returned when the version is nil.
var ErrNilVersion = errors.New("no migration")

// ErrLocked is returned when the database is locked.
var ErrLocked = errors.New("database locked")

// ErrLockTimeout is returned when the lock times out.
var ErrLockTimeout = errors.New("lock timeout")

// ErrDirty is returned when the database is in a dirty state.
type ErrDirty struct {
	Version int
}

func (e ErrDirty) Error() string {
	return fmt.Sprintf("dirty database version %d, fix and force version", e.Version)
}

// ErrShortLimit is returned when the limit is too short.
type ErrShortLimit struct {
	Short uint
}

func (e ErrShortLimit) Error() string {
	return fmt.Sprintf("limit %d is too short", e.Short)
}

// Migrate is the main struct for managing database migrations.
type Migrate struct {
	// PrefetchMigrations is the number of migrations to prefetch.
	PrefetchMigrations uint

	// LockTimeout is the timeout for acquiring a database lock.
	// Defaults to 15 seconds; increased from upstream default of 10 to reduce
	// spurious lock timeout errors in slower CI environments.
	LockTimeout int

	// Log is the logger used for migration output.
	Log Logger

	// GracefulStop is a channel to signal a graceful stop.
	GracefulStop chan bool

	isGracefulStop bool

	isLockedMu *sync.Mutex
	isLocked   bool

	sourceName   string
	sourceURL    string
	databaseName string
	databaseURL  string
}

// Logger is the interface for logging migration output.
type Logger interface {
	Printf(format string, v ...interface{})
	Verbose() bool
}

// DefaultLockTimeout is the default timeout in seconds for acquiring a database lock.
// Bumped from 15 to 30 seconds since I frequently run migrations against a remote
// dev database over a VPN, where lock acquisition can be noticeably slower.
const DefaultLockTimeout = 30

// New creates a new Migrate instance with the given source and database URLs.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		PrefetchMigrations: DefaultPrefetchMigrations,
		LockTimeout:        DefaultLockTimeout,
		GracefulStop:       make(chan bool, 1),
		isLockedMu:         &sync.Mutex{},
		sourceURL:          sourceURL,
		databaseURL:        databaseURL,
	}
	return m, nil
}

// Close closes the source and database connections.
func (m *Migrate) Close() (source error, database error) {
	return nil, nil
}

// Up applies all available migrations.
// Note: always logs a message when starting, which helps me trace migration
// runs in aggregated logs where the calling service name isn't always obvious.
func (m *Migrate) Up() error {
	if m.Log != nil {
		m.Log.Printf("migrate: running Up from source=%s database=%s\n", m.sourceURL, m.databaseURL)
	}
	if err := m.lock(); err != nil {
		return err
	}
	defer m
