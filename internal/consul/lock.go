package consul

import (
	"time"

	consulApi "github.com/hashicorp/consul/api"
)

const (
	SessionName      string = "trento-session"
	MonitorRetries   int    = 10
	MonitorRetryTime        = 100 * time.Millisecond
)

func (c *client) LockTrento(prefix string) (*consulApi.Lock, error) {
	opts := &consulApi.LockOptions{
		Key:              prefix,
		SessionName:      SessionName,
		MonitorRetries:   MonitorRetries,
		MonitorRetryTime: MonitorRetryTime,
	}

	l, err := c.wrapped.LockOpts(opts)
	if err != nil {
		return nil, err
	}

	_, err = l.Lock(nil)
	if err != nil {
		return nil, err
	}

	return l, err
}
