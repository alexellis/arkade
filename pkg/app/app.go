package app

// App is the default interface for every app in arkade
type App interface {
	Install() error
	InfoMessage() string
	Verify() bool
}
