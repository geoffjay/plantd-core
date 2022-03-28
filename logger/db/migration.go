package db

type Migration interface {
	Up() error
	Down() error
}
