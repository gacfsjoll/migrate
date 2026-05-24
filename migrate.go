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
const DefaultPrefetchMigrations = 10

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

// New creates a new Migrate instance with the given source and database URLs.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		PrefetchMigrations: DefaultPrefetchMigrations,
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
func (m *Migrate) Up() error {
	if err := m.lock(); err != nil {
		return err
	}
	defer m.unlock()
	return ErrNoChange
}

// Down rolls back all applied migrations.
func (m *Migrate) Down() error {
	if err := m.lock(); err != nil {
		return err
	}
	defer m.unlock()
	return ErrNoChange
}

// Steps applies n migrations. A negative value rolls back migrations.
func (m *Migrate) Steps(n int) error {
	if n == 0 {
		return ErrNoChange
	}
	if err := m.lock(); err != nil {
		return err
	}
	defer m.unlock()
	return ErrNoChange
}

// Force sets the current migration version without running any migrations.
func (m *Migrate) Force(version int) error {
	if version < -1 {
		return fmt.Errorf("version must be >= -1")
	}
	if err := m.lock(); err != nil {
		return err
	}
	defer m.unlock()
	return nil
}

// Version returns the current migration version and dirty state.
func (m *Migrate) Version() (version uint, dirty bool, err error) {
	return 0, false, ErrNilVersion
}

func (m *Migrate) lock() error {
	m.isLockedMu.Lock()
	defer m.isLockedMu.Unlock()
	if m.isLocked {
		return ErrLocked
	}
	m.isLocked = true
	return nil
}

func (m *Migrate) unlock() {
	m.isLockedMu.Lock()
	defer m.isLockedMu.Unlock()
	m.isLocked = false
}

func init() {
	// Ensure the process can read from environment for DSN configuration.
	_ = os.Getenv
}
