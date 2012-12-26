package pool

import (
	"testing"
)

func TestBasicOperation(t *testing.T) {
	p := make(Pool, 1)
	m := new(Member)

	p.Put(m)

	if mi, e := p.Get(0); mi != m || e != nil {
		t.Fail()
	}
}

func TestNoMembers(t *testing.T) {
	p := make(Pool, 1)

	if _, e := p.Get(0); e == nil {
		t.Fail()
	}
}

func TestTimeout(t *testing.T) {
	p := make(Pool, 1)
	m := new(Member)

	p.Put(m)

	m, e := p.Get(0)
	m, e = p.Get(0)

	if e != ErrTimeout {
		t.Fail()
	}
}

func TestLimit(t *testing.T) {
	p := make(Pool, 0)
	m := new(Member)

	e := p.Put(m)

	if e != ErrLimit {
		t.Fail()
	}
}

func BenchmarkOperation(b *testing.B) {
	b.StopTimer()
	p := make(Pool, b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		m := new(Member)
		e := p.Put(m)

		if e != nil {
			b.Fail()
		}
	}

	for i := 0; i < b.N; i++ {
		_, e := p.Get(0)

		if e != nil {
			b.Fail()
		}
	}
}
