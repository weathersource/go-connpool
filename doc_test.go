package connpool_test

import (
	"fmt"
	"sync"
	"time"

	connpool "github.com/weathersource/go-connpool"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type FooRequest struct{}
type FooResponse struct{}
type FooServiceClient interface {
	Foo(ctx context.Context, in *FooRequest, opts ...grpc.CallOption) (*FooResponse, error)
}
type fooServiceClient struct {
	cc *grpc.ClientConn
}

// faking a protobuf file import
type Pb struct{}

var pb = Pb{}

func (pb Pb) NewFooServiceClient(cc *grpc.ClientConn) FooServiceClient {
	return &fooServiceClient{cc}
}

func (c *fooServiceClient) Foo(ctx context.Context, in *FooRequest, opts ...grpc.CallOption) (*FooResponse, error) {
	return nil, nil
}

func Example() {
	var (
		wg          sync.WaitGroup
		reqCnt      = 100 // how many GRPC requests to make
		connFactory = func() (*grpc.ClientConn, error) {
			return grpc.Dial("localhost:50051", grpc.WithInsecure())
		}
		poolCapacity    = 3
		poolIdleTimeout = 2 * time.Second
		poolMaxLife     = 30 * time.Second
	)

	// Set up the response channel for goroutines
	type Res struct {
		out *FooResponse // the GRPC response message
		err error
	}
	c := make(chan Res, reqCnt)

	// configure the connection pool
	pool := connpool.New(
		connFactory,
		poolCapacity,
		poolIdleTimeout,
		poolMaxLife,
	)

	for i := 0; i < reqCnt; i++ {

		wg.Add(1)

		go func() {

			conn, err := pool.Get(context.Background())
			if nil != err {
				c <- Res{err: err}
			} else {
				client := pb.NewFooServiceClient(conn.ClientConn)
				out, err := client.Foo(context.Background(), &FooRequest{})
				c <- Res{
					out: out,
					err: err,
				}
			}

			conn.Close()
			wg.Done()
		}()
	}

	// close the channel once all request go routines are complete
	go func() {
		wg.Wait()
		close(c)
		pool.Close()
	}()

	for res := range c {
		if res.err != nil {
			fmt.Println("Error:", res.err)
		} else if res.out != nil {
			fmt.Println("Response:", res.out)
		}
	}
}
