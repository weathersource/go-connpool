package foo

import (
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

func NewFooServiceClient(cc *grpc.ClientConn) FooServiceClient {
	return &fooServiceClient{cc}
}

func (c *fooServiceClient) Foo(ctx context.Context, in *FooRequest, opts ...grpc.CallOption) (*FooResponse, error) {
	return nil, nil
}
