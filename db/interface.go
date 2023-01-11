package db

type DB interface {
	Test() (string, error)
}
