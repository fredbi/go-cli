package wait

import "time"

type (
	Option func(*options)

	options struct {
		timeout         time.Duration
		pollingInterval time.Duration
		dialProtocol    string
	}
)

func defaultOptions() *options {
	return &options{
		timeout:         5 * time.Second,
		pollingInterval: 100 * time.Millisecond,
		dialProtocol:    "tcp",
	}
}

func applyOptions(opts []Option) *options {
	o := defaultOptions()

	for _, apply := range opts {
		apply(o)
	}

	return o
}
