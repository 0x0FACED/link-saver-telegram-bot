package grpc

import (
	"context"

	"github.com/0x0FACED/proto-files/link_service/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type APIClient struct {
	client gen.LinkServiceClient
}

func New(host string) (*APIClient, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &APIClient{
		client: gen.NewLinkServiceClient(conn),
	}, nil
}

func (a APIClient) GetLink(ctx context.Context, req *gen.GetLinkRequest) (*gen.GetLinkResponse, error) {
	resp, err := a.client.GetLink(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a APIClient) GetLinks(ctx context.Context, req *gen.GetLinksRequest) (*gen.GetLinksResponse, error) {
	resp, err := a.client.GetLinks(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a APIClient) GetAllLinks(ctx context.Context, req *gen.GetAllLinksRequest) (*gen.GetAllLinksResponse, error) {
	resp, err := a.client.GetAllLinks(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a APIClient) SaveLink(ctx context.Context, req *gen.SaveLinkRequest) (*gen.SaveLinkResponse, error) {
	resp, err := a.client.SaveLink(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a APIClient) DeleteLink(ctx context.Context, req *gen.DeleteLinkRequest) (*gen.DeleteLinkResponse, error) {
	resp, err := a.client.DeleteLink(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
