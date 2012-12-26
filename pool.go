package pool

import (
	"errors"
	"time"
)

type Pool struct {
	Size    int
	channel chan interface{}
}

var ErrNoMember = errors.New("pool: could not obtain member")
var ErrLimit = errors.New("pool: already at size limit")
var ErrTimeout = errors.New("pool: retrieval timed out")

// 
func NewPool(max int) *Pool {
	return &Pool{Size: 0, channel: make(chan interface{}, max)}
}

func (p *Pool) Get(timeoutDuration time.Duration) (interface{}, error) {
	if len(p.channel) == 0 && p.Size < cap(p.channel) {
		return nil, ErrNoMember
	}

	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(timeoutDuration)
		timeout <- true
	}()

	select {
	case m := <-p.channel:
		return m, nil
	case <-timeout:
		return nil, ErrTimeout
	}

	return nil, nil
}

func (p *Pool) Register(m interface{}) error {
	if p.Size > cap(p.channel) || cap(p.channel) == 0 {
		return ErrLimit
	}

	p.Size = p.Size + 1
	return nil
}

func (p *Pool) Put(m interface{}) error {
	if e := p.Register(m); e != nil {
		return e
	}

	p.channel <- m

	return nil
}
