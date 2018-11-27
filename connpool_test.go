package connpool

import (
	"context"
	"errors"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
	grpc "google.golang.org/grpc"
)

// Test 0:
// Create a pool with 0 capacity, ensure it is set to 1 capaicity
// Close an empty pool
// Call Get on an empty pool
func TestPool0(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	factory := func() (*grpc.ClientConn, error) { return nil, nil }
	pool := New(
		factory,
		0,
		400*time.Microsecond,
		800*time.Microsecond,
	)
	assert.Equal(0, len(pool.conns))
	assert.Equal(1, cap(pool.conns))

	// close an empty pool
	pool.Close()
	assert.True(pool.isClosed())

	// call Get on a closed pool
	conn, err := pool.Get(ctx)
	assert.Nil(conn)
	assert.NotNil(err)
}

// Test 1:
// Create a pool with a factory that returns an error
// Call get that returns an error
// Close the pool
func TestPool1(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// test factory that returns error
	factory := func() (*grpc.ClientConn, error) { return nil, errors.New("foo") }
	pool := New(
		factory,
		2,
		400*time.Microsecond,
		800*time.Microsecond,
	)
	assert.Equal(0, len(pool.conns))
	assert.Equal(2, cap(pool.conns))

	// Call get that returns an error
	conn, err := pool.Get(ctx)
	assert.Nil(conn)
	assert.NotNil(err)

	// close an empty pool
	pool.Close()
	assert.True(pool.isClosed())
}

// Test 2:
// Create lots of connections to test quasi-load-balancing code
// Close all clients for a healthy connection
// Create unhealthy connection with no clients
// Create unhealthy connection with clients
// Close all clients for a quarintined connection
// Close pool with non-empty healthy connection
// Close pool with non-empty quarintine connection
func TestPool2(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	factory := func() (*grpc.ClientConn, error) { return nil, nil }
	pool := New(
		factory,
		4,
		500*time.Microsecond,
		1000*time.Microsecond,
	)
	assert.Equal(0, len(pool.conns))
	assert.Equal(4, cap(pool.conns))

	// Get 9 connections to test quasi-load-balancing code
	conn0, err0 := pool.Get(ctx)
	assert.NotNil(conn0)
	assert.Nil(err0)
	conn1, err1 := pool.Get(ctx)
	assert.NotNil(conn1)
	assert.Nil(err1)
	conn2, err2 := pool.Get(ctx)
	assert.NotNil(conn2)
	assert.Nil(err2)
	conn3, err3 := pool.Get(ctx)
	assert.NotNil(conn3)
	assert.Nil(err3)
	conn4, err4 := pool.Get(ctx)
	assert.NotNil(conn4)
	assert.Nil(err4)
	conn5, err5 := pool.Get(ctx)
	assert.NotNil(conn5)
	assert.Nil(err5)
	conn6, err6 := pool.Get(ctx)
	assert.NotNil(conn6)
	assert.Nil(err6)
	conn7, err7 := pool.Get(ctx)
	assert.NotNil(conn7)
	assert.Nil(err7)
	conn8, err8 := pool.Get(ctx)
	assert.NotNil(conn8)
	assert.Nil(err8)

	// There are now 4 connections in the pool, that are shared accordingly:
	// - Group 0: conn0, conn4, conn8
	// - Group 1: conn1, conn5
	// - Group 2: conn2, conn6
	// - Group 3: conn3, conn7
	assert.Equal(4, len(pool.conns))
	assert.Equal(3, conn0.clientCount)
	assert.Equal(2, conn1.clientCount)
	assert.Equal(2, conn2.clientCount)
	assert.Equal(2, conn3.clientCount)
	assert.Equal(0, len(pool.quarintine))

	// Close all clients for a healthy connection (Group 1)
	conn1.Close()
	assert.Equal(1, conn1.clientCount)
	conn5.Close()
	assert.Equal(0, conn5.clientCount)

	// Sleep to cause unhealthy conns
	time.Sleep(750 * time.Microsecond)

	// Create unhealthy connection with no clients
	// This will send Group 1 to quarintine
	// but since it is empty, it will just get closed
	conn9, err9 := pool.Get(ctx)
	assert.NotNil(conn9)
	assert.Nil(err9)
	assert.Equal(0, len(pool.quarintine))

	// Create unhealthy connection with clients
	// This will send Group 2 to quarintine
	conn10, err10 := pool.Get(ctx)
	assert.NotNil(conn10)
	assert.Nil(err10)
	assert.Equal(1, len(pool.quarintine))

	// Close all clients for a quarintined connection (Group 2)
	// This will close out Group 2 and delete it from quarintine
	conn2.Close()
	assert.Equal(1, conn2.clientCount)
	conn6.Close()
	assert.Equal(0, conn6.clientCount)
	assert.Equal(0, len(pool.quarintine))

	// Create unhealthy connection with clients
	// This will send Group 3 to quarintine
	conn11, err11 := pool.Get(ctx)
	assert.NotNil(conn11)
	assert.Nil(err11)
	assert.Equal(1, len(pool.quarintine))

	// We still have valid connections in Group 1, plus conn9, conn10, and conn11
	// We still have connections in quarintine (Group 3)
	// Close the pool
	pool.Close()
	assert.True(pool.isClosed())
}
