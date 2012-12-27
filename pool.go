// Package pool provides a basic pool mechanism handy for tasks like
// database connection management.
package pool

import (
	"errors"
	"sync"
	"time"
)

type Pool struct {
	// Size represents the current number of resources checked into 
	// the pool, including resources that are checked out. 
	Size int

	// channel is the mechanism by which the pool stores and awaits
	// available resources
	channel chan interface{}
	mutex   sync.Mutex
}

// ErrNoMember is returned on a Get operation when there are no 
// available resources to check out from the pool. If this error
// is returned the calling procedure may create a new pool resource
// and check it in.
var ErrNoMember = errors.New("pool: could not obtain member")

// ErrLimit indicates that the pool is already at the resource
// limit and cannot check in new values. 
var ErrLimit = errors.New("pool: already at size limit")

// ErrTimeout is returned on a Get operation when receiving an 
// item from the pool has timed out.
var ErrTimeout = errors.New("pool: retrieval timed out")

// New returns a new Pool with the specified maximum
// size of stored resources. 
func New(max int) *Pool {
	return &Pool{Size: 0, channel: make(chan interface{}, max)}
}

// Get returns a new connection from the pool. The timeoutDuration
// specifies the time before attempted retrieval of the connection
// fails with ErrTimeout.
func (p *Pool) Get(timeoutDuration time.Duration) (interface{}, error) {
	p.mutex.Lock()

	if len(p.channel) == 0 && p.Size < cap(p.channel) {
		p.mutex.Unlock()
		return nil, ErrNoMember
	}

	p.mutex.Unlock()

	select {
	case m := <-p.channel:
		return m, nil
	case <-time.After(timeoutDuration):
		return nil, ErrTimeout
	}

	return nil, nil
}

// Register provides a mechanism for adding a resource to the pool,
// but ensuring it is not available for checkout during that time. 
// It should be used whenever the added resource needs to act as 
// though it's been checked out already.
func (p *Pool) Register(m interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Size > cap(p.channel) || cap(p.channel) == 0 {
		return ErrLimit
	}

	p.Size = p.Size + 1
	return nil
}

// Registers a new resource and adds it to the pool for checkout
// with Get later.
func (p *Pool) Put(m interface{}) error {
	if e := p.Register(m); e != nil {
		return e
	}

	p.channel <- m

	return nil
}
