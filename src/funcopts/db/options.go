package db

import "log"

type Option interface {
	apply(*config) error
}

type option func(*config) error

func (f option) apply(cfg *config) error {
	return f(cfg)
}

type config struct {
	cache  bool
	logger *log.Logger
}

func WithCache() Option {
	f := func(cfg *config) error {
		cfg.cache = true
		return nil
	}
	return option(f)
}

func WithLogger(logger *log.Logger) Option {
	f := func(cfg *config) error {
		cfg.logger = logger
		return nil
	}
	return option(f)
}
