package connpool

import (
	"context"
	"math/rand"
	"sync"
	"time"

	errors "github.com/weathersource/go-errors"
	grpc "google.golang.org/grpc"
)

// Pool is the GRPC connection pool.
type Pool struct {
	sync.RWMutex
	conns       chan *ClientConn
	quarintine  map[int]*ClientConn
	factory     Factory
	idleTimeout time.Duration
	lifeTimeout time.Duration
	closed      bool
}

// Factory is a function type creating a GRPC connection.
type Factory func() (*grpc.ClientConn, error)

// New creates a new connection pool. The factory parameter accpets a function
// that creates new connections. The capacity parameter sets the maximum
// capacity of healthy connections inside the connection pool. The idleTimeout
// parameter sets a duration after which, if no new requests are made on the
// connection, the connection is recycled. The lifeTimeout parameter sets a
// duration that is the maximum lifetime of a connection before it is recycled.
func New(
	factory Factory,
	capacity int,
	idleTimeout time.Duration,
	lifeTimeout time.Duration) *Pool {

	if capacity <= 0 {
		capacity = 1
	}
	p := &Pool{
		conns:       make(chan *ClientConn, capacity),
		quarintine:  make(map[int]*ClientConn),
		factory:     factory,
		idleTimeout: idleTimeout,
		lifeTimeout: lifeTimeout,
	}
	return p
}

// Get retreives a healthy connection for use with a GRPC client.
func (p *Pool) Get(context.Context) (*ClientConn, error) {
	var (
		err            error
		conn           *ClientConn
		now            = time.Now()
		jitter         = rand.Float64()*.25 + .75
		expireTime     = now.Add(time.Duration(float64(p.lifeTimeout) * jitter))
		capacity       = cap(p.conns)
		avgClientCount int
	)

	p.Lock()
	defer p.Unlock()

	// if channel is closed, return error
	if p.isClosed() {
		return nil, errors.NewInternalError("Connection pool is closed.")
	}

	// reuse a connection if the pool is full
	if len(p.conns) == capacity {

		// this adds not-entirely-dumb load balancing to connection selection
		for i := 0; i < capacity+1; i++ {
			conn = <-p.conns

			if conn.checkHealth() {
				p.conns <- conn
				conn.RLock()
				clientCount := conn.clientCount
				conn.RUnlock()
				if clientCount < avgClientCount ||
					clientCount <= 1 ||
					i == capacity {
					conn.addClient()
					break
				}
				avgClientCount = (avgClientCount*i + clientCount) / (i + 1)
			} else {
				conn.quarintine()
				break
			}
		}
	}

	// If there is unused capacity in the connection pool, add a new connection.
	// This is purposefully not an else to the previous if statement, as it is
	// possible that we could enter that if statement with full capacity, and
	// an unhealthy connection is removed. We would then reach this statement
	// with a non-full connection pool.
	if len(p.conns) < capacity {
		conn = &ClientConn{
			pool:        p,
			expireTime:  expireTime,
			lastTime:    now,
			clientCount: 1,
		}
		conn.ClientConn, err = p.factory()
		if err != nil {
			return nil, err
		}
		p.conns <- conn
	}

	return conn, nil
}

// Close closes the connection pool, and terminates all connections contained
// within the connection pool.
func (p *Pool) Close() {

	// The zero value of a channel is nil, so conns is initally nil
	var (
		conns chan *ClientConn
		q     map[int]*ClientConn
	)

	p.Lock()
	if p.closed == false {
		conns = p.conns
		p.conns = nil
		q = p.quarintine
		p.quarintine = nil
		p.closed = true
	}
	p.Unlock()

	// If conns is nil, then the pool was already closed.
	if conns != nil {
		// Close the channel and shut down all the connections within it.
		close(conns)
		for conn := range conns {
			if conn.ClientConn != nil {
				conn.ClientConn.Close()
			}
		}

		// close all conns in quarintine
		for _, conn := range q {
			if conn.ClientConn != nil {
				conn.ClientConn.Close()
			}
		}
	}
}

// isClosed returns if a connection pool is closed.
func (p *Pool) isClosed() bool {
	return p.closed
}

// ClientConn is a client connection managed by the connection pool.
type ClientConn struct {
	*grpc.ClientConn
	sync.RWMutex
	pool          *Pool
	quarintineKey int
	expireTime    time.Time
	lastTime      time.Time
	clientCount   int
}

// Close communicates to a connection that a GRPC client is done using the
// connection. The connection itself will only be terminated if it is no longer
// healthy.
func (c *ClientConn) Close() {
	c.Lock()
	if c.clientCount > 0 {
		c.clientCount--
		if c.quarintineKey != 0 && c.clientCount == 0 {
			if c.ClientConn != nil {
				c.ClientConn.Close()
			}
			c.pool.Lock()
			delete(c.pool.quarintine, c.quarintineKey)
			c.pool.Unlock()
		}
	}
	c.Unlock()
}

// addClient communicates to a connection that a new GRPC client is utelizing
// the connection.
func (c *ClientConn) addClient() {
	c.Lock()
	defer c.Unlock()
	c.clientCount++
	c.lastTime = time.Now()
}

// checkHealth communicates if a connection is healthy or not.
func (c *ClientConn) checkHealth() bool {
	c.RLock()
	defer c.RUnlock()
	now := time.Now()
	return c.expireTime.After(now) && c.lastTime.Add(c.pool.idleTimeout).After(now)
}

// quarintine removes a connection from the connection pool, and places it in
// quarintine where existing clients can complete the requests they have made on
// that connection. After all clients are done with the connection, it will be
// terminated.
func (c *ClientConn) quarintine() {
	c.Lock()
	defer c.Unlock()
	if c.clientCount == 0 {
		if c.ClientConn != nil {
			c.ClientConn.Close()
		}
	} else {
		for {
			key := rand.Int()
			_, ok := c.pool.quarintine[key]
			if !ok && key > 0 {
				c.pool.quarintine[key] = c
				c.quarintineKey = key
				break
			}
		}
	}
}
