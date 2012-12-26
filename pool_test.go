package pool

import (
	"fmt"
	"testing"
	"time"
)

type MemberObject struct {
	id int
}

func TestBasicOperation(t *testing.T) {
	p := NewPool(1)
	m := &MemberObject{id: 1}

	p.Put(m)

	mi, e := p.Get(0)

	if e != nil {
		t.Fail()
	}

	if mi.(*MemberObject).id != 1 {
		t.Fail()
	}
}

func TestRegister(t *testing.T) {
	p := NewPool(1)
	m := &MemberObject{id: 1}

	p.Register(m)

	if p.Size != 1 {
		t.Fail()
	}
}

func TestNoMembers(t *testing.T) {
	p := NewPool(1)

	if _, e := p.Get(0); e == nil {
		t.Fail()
	}
}

func TestTimeout(t *testing.T) {
	p := NewPool(1)
	m := &MemberObject{id: 1}

	p.Put(m)

	var e error

	_, e = p.Get(0)
	_, e = p.Get(0)

	if e != ErrTimeout {
		t.Fail()
	}
}

func TestLimit(t *testing.T) {
	p := NewPool(0)
	m := &MemberObject{id: 1}

	e := p.Put(m)

	if e != ErrLimit {
		t.Fail()
	}
}

func BenchmarkOperation(b *testing.B) {
	b.StopTimer()
	p := NewPool(b.N)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		m := &MemberObject{id: 1}
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

func ExamplePool() {
	type Connection struct{}

	// Generate a new connection
	createConnection := func() *Connection {
		fmt.Println("New connection")
		return &Connection{}
	}

	// Our pool has a maximum of 2 connections
	pool := NewPool(2)

	// Grab a connection from the pool, or try to add a new connection
	getConnection := func() (*Connection, error) {
		m, e := pool.Get(0 * time.Second)

		if m == nil {
			if e == ErrNoMember {
				c := createConnection()
				pool.Register(c)
				return c, nil

			} else {
				fmt.Println("Timed out")
				return nil, e
			}
		}

		fmt.Println("Using existing connection")
		return m.(*Connection), nil
	}

	// Check out 2 new connections
	c1, _ := getConnection()
	getConnection()

	// Put one back in
	pool.Put(c1)

	// Check the one we put back in out
	c1, _ = getConnection()

	// Time out waiting for another
	getConnection()

	// Output:
	// New connection
	// New connection
	// Using existing connection
	// Timed out
}
