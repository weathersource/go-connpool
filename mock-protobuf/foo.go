package foo

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// FooRequest is a mock protobuf request message
type FooRequest struct{}

// FooRequest is a mock protobuf response message
type FooResponse struct{}

// FooServiceClient is a mock protobuf service client interface
type FooServiceClient interface {
	Foo(ctx context.Context, in *FooRequest, opts ...grpc.CallOption) (*FooResponse, error)
}

type fooServiceClient struct {
	cc *grpc.ClientConn
}

// NewFooServiceClient creates a mock protobuf service client
func NewFooServiceClient(cc *grpc.ClientConn) FooServiceClient {
	return &fooServiceClient{cc}
}

// Foo is a mock protobuf RPC
func (c *fooServiceClient) Foo(ctx context.Context, in *FooRequest, opts ...grpc.CallOption) (*FooResponse, error) {
	return nil, nil
}
