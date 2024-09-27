package postgres

import "time"

// Option -.
type Option func(*Postgres)

func OptionSet(size, attempts int, timeout time.Duration) Option {
	return func(c *Postgres) {
		MaxPoolSize(size)(c)
		ConnAttempts(attempts)(c)
		ConnTimeout(timeout)(c)
	}
}

// MaxPoolSize Максимальный размер пула
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts Попытки соединения
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout Время ожидания соединения
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
