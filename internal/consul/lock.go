package consul

import (
	"time"

	consulApi "github.com/hashicorp/consul/api"
)

const (
	sessionName      string = "trento-session"
	monitorRetries   int    = 10
	monitorRetryTime        = 100 * time.Millisecond
)

func (c *client) LockTrento(prefix string) (*consulApi.Lock, error) {
	opts := &consulApi.LockOptions{
		Key:              prefix,
		SessionName:      sessionName,
		MonitorRetries:   monitorRetries,
		MonitorRetryTime: monitorRetryTime,
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
