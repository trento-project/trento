package consul

import (
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var fakeLockGetAttemps int = 0
var globalT *testing.T

func increaseAttempts() {
	fakeLockGetAttemps++
}

func fakeLockGet(c *client, prefix string, q *consulApi.QueryOptions) (*consulApi.KVPair, *consulApi.QueryMeta, error) {
	defer increaseAttempts()

	assert.Equal(globalT, prefix, "test")

	if fakeLockGetAttemps == -2 {
		kvPair := &consulApi.KVPair{
			Flags:   1,
			Session: "session",
		}
		qMeta := &consulApi.QueryMeta{}

		return kvPair, qMeta, nil
	} else if fakeLockGetAttemps == -1 {
		kvPair := &consulApi.KVPair{}
		qMeta := &consulApi.QueryMeta{}

		return kvPair, qMeta, errors.New("error")
	} else if fakeLockGetAttemps == 0 {
		kvPair := &consulApi.KVPair{
			Flags:   consulApi.LockFlagValue,
			Session: "session",
		}
		qMeta := &consulApi.QueryMeta{
			LastIndex: 10,
		}

		return kvPair, qMeta, nil
	} else if fakeLockGetAttemps == 1 {
		kvPair := &consulApi.KVPair{
			Flags:   consulApi.LockFlagValue,
			Session: "session",
		}
		qMeta := &consulApi.QueryMeta{
			LastIndex: 11,
		}

		assert.Equal(globalT, uint64(0xa), q.WaitIndex)

		return kvPair, qMeta, nil
	} else {
		kvPair := &consulApi.KVPair{
			Flags:   consulApi.LockFlagValue,
			Session: "",
		}
		qMeta := &consulApi.QueryMeta{}

		assert.Equal(globalT, uint64(0xb), q.WaitIndex)

		return kvPair, qMeta, nil
	}
}

func TestLockWaitReleased(t *testing.T) {
	globalT = t
	fakeLockGetAttemps = 0
	c, _ := DefaultClient()
	lockGet = fakeLockGet

	err := c.LockWaitReleasead("test")

	assert.Equal(t, fakeLockGetAttemps, 3)
	assert.NoError(t, err)
}

func TestLockWaitReleasedErr(t *testing.T) {
	globalT = t
	fakeLockGetAttemps = -1
	c, _ := DefaultClient()
	lockGet = fakeLockGet

	err := c.LockWaitReleasead("test")

	assert.Equal(t, fakeLockGetAttemps, 0)
	assert.EqualError(t, err, "failed to read lock: error")
}

func TestLockWaitReleasedLockConflict(t *testing.T) {
	globalT = t
	fakeLockGetAttemps = -2
	c, _ := DefaultClient()
	lockGet = fakeLockGet

	err := c.LockWaitReleasead("test")

	assert.Equal(t, fakeLockGetAttemps, -1)
	assert.EqualError(t, err, "Existing key does not match lock use")
}

func fakeLock(l *consulApi.Lock, stopCh <-chan struct{}) (<-chan struct{}, error) {
	return nil, nil
}

func fakeLockErr(l *consulApi.Lock, stopCh <-chan struct{}) (<-chan struct{}, error) {
	return nil, errors.New("error")
}

func TestLockTrento(t *testing.T) {
	globalT = t
	c, _ := DefaultClient()
	lock = fakeLock

	l, err := c.LockTrento("test")

	assert.IsType(t, &consulApi.Lock{}, l)
	assert.NoError(t, err)
}

func TestLockTrentoNoKey(t *testing.T) {
	globalT = t
	c, _ := DefaultClient()
	lock = fakeLock

	l, err := c.LockTrento("")

	assert.Equal(t, (*consulApi.Lock)(nil), l)
	assert.EqualError(t, err, "missing key")
}

func TestLockTrentoErr(t *testing.T) {
	globalT = t
	c, _ := DefaultClient()
	lock = fakeLockErr

	l, err := c.LockTrento("test")

	assert.Equal(t, (*consulApi.Lock)(nil), l)
	assert.EqualError(t, err, "error")
}
