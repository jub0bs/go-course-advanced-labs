package db

import "log"

type Connection struct {
	// omitted fields
}

func (*Connection) String() string {
	return "Connection established!"
}

func Open(addr string, opts ...Option) (*Connection, error) {
	var cfg config
	for _, o := range opts {
		err := o.applyTo(&cfg)
		if err != nil {
			return nil, err
		}
	}
	var conn Connection
	return &conn, nil
}

type config struct {
	cache  bool
	logger log.Logger
}
