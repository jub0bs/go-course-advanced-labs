package db

import "log"

type Option interface {
	applyTo(*config) error
}

type cacheOption struct{}

func (c *cacheOption) applyTo(cfg *config) error {
	cfg.cache = true
	return nil
}

func WithCache() Option {
	return new(cacheOption)
}

type loggerOption log.Logger

func (l loggerOption) applyTo(cfg *config) error {
	cfg.logger = log.Logger(l)
	return nil
}

func WithLogger(log log.Logger) Option {
	return loggerOption(log)
}
