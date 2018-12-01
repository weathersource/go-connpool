package connpool_test

import (
	"fmt"
	"sync"
	"time"

	connpool "github.com/weathersource/go-connpool"
	pb "github.com/weathersource/go-connpool/foo-proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

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

	// set up the response channel for goroutines
	type Res struct {
		out *pb.BarResponse // the GRPC response message
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

	// concurrent GRPC requests
	for i := 0; i < reqCnt; i++ {
		wg.Add(1)
		go func() {
			conn, err := pool.Get(context.Background())
			if nil != err {
				c <- Res{err: err}
			} else {
				client := pb.NewFooClient(conn.ClientConn)
				out, err := client.Bar(context.Background(), &pb.BarRequest{})
				c <- Res{
					out: out,
					err: err,
				}
			}
			conn.Close()
			wg.Done()
		}()
	}

	// close the channel and connection pool once all GRPC requests are complete
	go func() {
		wg.Wait()
		close(c)
		pool.Close()
	}()

	// handle GRPC responses
	for res := range c {
		if res.err != nil {
			fmt.Println("Error:", res.err)
		} else if res.out != nil {
			fmt.Println("Response:", res.out)
		}
	}
}
