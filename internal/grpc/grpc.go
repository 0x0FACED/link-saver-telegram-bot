package grpc

import (
	"context"

	"github.com/0x0FACED/proto-files/link_service/gen"
	"google.golang.org/grpc"
)

type APIClient struct {
	client gen.LinkServiceClient
}

func New(conn *grpc.ClientConn) *APIClient {
	return &APIClient{
		client: gen.NewLinkServiceClient(conn),
	}
}

func (a APIClient) GetLinks(ctx context.Context, req gen.GetLinksRequest) (*gen.GetLinksResponse, error) {
	panic("TODO: impl me")
}

func (a APIClient) SaveLink(ctx context.Context, req gen.SaveLinkRequest) (*gen.SaveLinkResponse, error) {
	panic("TODO: impl me")
}

func (a APIClient) DeleteLink(ctx context.Context, req gen.DeleteLinkRequest) (*gen.DeleteLinkResponse, error) {
	panic("TODO: impl me")
}
