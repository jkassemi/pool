package pool

import (
  "testing"
)

func TestBasicOperation(t *testing.T){
  p := make(Pool, 1)
  m := new(Member)

  p.Put(m)

  if mi, e := p.Get(0); mi != m || e != nil {
    t.Failed()
  }
}

func TestNoMembers(t *testing.T){
  p := make(Pool, 1)

  if _, e := p.Get(0); e == nil {
    t.Failed()
  }
}

func TestTimeout(t *testing.T){
  p := make(Pool, 1)

  _, e := p.Get(0)
  _, e = p.Get(0)

  if e != ErrTimeout {
    t.Failed()
  }
}

func TestLimit(t *testing.T){
  p := make(Pool, 0)
  m := new(Member)

  e := p.Put(m)

  if e != ErrLimit {
    t.Failed()
  }
}
