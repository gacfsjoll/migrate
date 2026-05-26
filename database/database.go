// Package database provides the interface and registry for database drivers
// used by the migrate tool to apply migrations.
package database

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

// NilVersion is the version when no migration has been applied yet.
const NilVersion = -1

// Driver is the interface that database drivers must implement to be
// compatible with the migrate tool.
type Driver interface {
	// Open returns a new driver instance configured with the given URL.
	// The URL format is driver-specific.
	Open(url string) (Driver, error)

	// Close releases any resources held by the driver.
	Close() error

	// Lock acquires an exclusive lock on the database to prevent concurrent
	// migrations from running simultaneously.
	Lock() error

	// Unlock releases the exclusive lock acquired by Lock.
	Unlock() error

	// Run applies the given migration reader to the database.
	Run(migration io.Reader) error

	// SetVersion stores the current migration version and dirty state
	// in the database.
	SetVersion(version int, dirty bool) error

	// Version returns the currently applied migration version and whether
	// the database is in a dirty state. Returns NilVersion if no migration
	// has been applied.
	Version() (version int, dirty bool, err error)

	// Drop deletes all database resources managed by the driver, effectively
	// resetting the database to its initial state.
	Drop() error
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("database: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("database: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Open returns a new database driver instance for the given URL.
// The driver name is extracted from the URL scheme.
func Open(url string) (Driver, error) {
	scheme, err := schemeFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("database: failed to parse URL scheme: %w", err)
	}

	driversMu.RLock()
	d, ok := drivers[scheme]
	driversMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("database: unknown driver %q (forgotten import?)", scheme)
	}

	return d.Open(url)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()

	list := make([]string, 0, len(drivers))
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// schemeFromURL extracts the scheme (driver name) from a database URL.
func schemeFromURL(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	for i := 0; i < len(url); i++ {
		switch url[i] {
		case ':':
			if i == 0 {
				return "", fmt.Errorf("missing scheme in URL")
			}
			return url[:i], nil
		case '/', '?', '#':
			return "", fmt.Errorf("missing scheme in URL")
		}
	}

	return "", fmt.Errorf("missing scheme in URL")
}
