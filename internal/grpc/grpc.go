package grpc

import (
	"context"

	"github.com/0x0FACED/link-saver-telegram-bot/config"
	pdf "github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type APIClient struct {
	client gen.LinkServiceClient

	pdfClient pdf.PDFServiceClient
}

func New(cfg config.Config) (*APIClient, error) {
	linkService, err := linkService(cfg.GRPC.Link.Host, cfg.GRPC.Link.Port)
	if err != nil {
		return nil, err
	}

	pdfService, err := pdfService(cfg.GRPC.PDF.Host, cfg.GRPC.PDF.Port)
	if err != nil {
		return nil, err
	}

	return &APIClient{
		client:    linkService,
		pdfClient: pdfService,
	}, nil
}

func linkService(host string, port string) (gen.LinkServiceClient, error) {
	addr := host + ":" + port
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return gen.NewLinkServiceClient(conn), nil
}

const maxMsgSize = 15 * 1024 * 1024 // 15mb

func pdfService(host string, port string) (pdf.PDFServiceClient, error) {
	addr := host + ":" + port
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
	)

	if err != nil {
		return nil, err
	}

	return pdf.NewPDFServiceClient(conn), nil
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

// PDF ====================================================

func (a APIClient) ConvertToPDF(ctx context.Context, req *pdf.ConvertToPDFRequest) (*pdf.ConvertToPDFResponse, error) {
	resp, err := a.pdfClient.ConvertToPDF(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (a APIClient) GetSavedPDF(ctx context.Context, req *pdf.GetSavedPDFRequest) (*pdf.GetSavedPDFResponse, error) {
	resp, err := a.pdfClient.GetSavedPDF(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (a APIClient) DeletePDF(ctx context.Context, req *pdf.DeletePDFRequest) (*pdf.DeletePDFResponse, error) {
	resp, err := a.pdfClient.DeletePDF(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (a APIClient) DeleteAllPDF(ctx context.Context, req *pdf.DeleteAllPDFRequest) (*pdf.DeleteAllPDFResponse, error) {
	resp, err := a.pdfClient.DeleteAllPDF(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
