package databasex

type Config struct {
	Debug        bool
	DSN          string
	MaxOpenConns uint
	MaxIdleConns uint
}
