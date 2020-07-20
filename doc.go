// Package synchrozine provides an instrument for synchronization of multiple goroutines over a single channel.
// It provides the main channel (`chan error`), as well as tools for complete synchronization and receivers a channels list to send finish signals to goroutines.
//
// Synchrozine supports the startup synchronization and thread-safe injections.
package synchrozine
