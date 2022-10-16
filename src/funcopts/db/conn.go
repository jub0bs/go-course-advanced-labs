package db

type Connection struct {
	// simplistic connection for the sake of this exercise
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
