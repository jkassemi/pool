package pool

import (
  "errors"
  "time"
)

type Pool chan *Member
type Member interface{}

var ErrNoMembers = errors.New("pool: could not obtain member")
var ErrLimit = errors.New("pool: already at size limit")
var ErrTimeout = errors.New("pool: retrieval timed out")

func (p Pool) Get(timeoutDuration time.Duration) (*Member, error) {
  if len(p) == 0 {
    return nil, ErrNoMembers
  }

  timeout := make(chan bool, 1)

  go func() {
    time.Sleep(timeoutDuration)
    timeout <- true
  }()

  select {
    case m := <-p:
      return m, nil
    case <-timeout:
      return nil, ErrTimeout
  }

  return nil, nil
}


func (p Pool) Put(m *Member) (error) {
  if len(p) == cap(p) {
    return ErrLimit
  }

  p <- m

  return nil
}
