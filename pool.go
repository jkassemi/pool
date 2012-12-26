package pool

import (
	"errors"
	"time"
)

type Pool struct {
	Size    int
	Channel chan interface{}
}

var ErrNoMember = errors.New("pool: could not obtain member")
var ErrLimit = errors.New("pool: already at size limit")
var ErrTimeout = errors.New("pool: retrieval timed out")

func NewPool(max int) *Pool {
	return &Pool{Size: 0, Channel: make(chan interface{}, max)}
}

func (p *Pool) Get(timeoutDuration time.Duration) (interface{}, error) {
	if len(p.Channel) == 0 && p.Size < cap(p.Channel) {
		return nil, ErrNoMember
	}

	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(timeoutDuration)
		timeout <- true
	}()

	select {
	case m := <-p.Channel:
		return m, nil
	case <-timeout:
		return nil, ErrTimeout
	}

	return nil, nil
}

func (p *Pool) Register(m interface{}) error {
	if p.Size > cap(p.Channel) || cap(p.Channel) == 0 {
		return ErrLimit
	}

	p.Size = p.Size + 1
	return nil
}

func (p *Pool) Put(m interface{}) error {
	if e := p.Register(m); e != nil {
		return e
	}

	p.Channel <- m

	return nil
}
