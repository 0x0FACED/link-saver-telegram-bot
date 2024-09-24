package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/grpc/mocks"
	"github.com/0x0FACED/pdf-proto/pdf_service/gen"
	pdf "github.com/0x0FACED/pdf-proto/pdf_service/gen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestConvertPDF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPDFService := mocks.NewMockPDFServiceClient(ctrl)

	request := &gen.ConvertToPDFRequest{
		UserId:      123,
		OriginalUrl: "http://example.com",
		Description: "test",
		Scale:       1.0,
	}
	expectedResponse := &gen.ConvertToPDFResponse{PdfData: []byte("pdf")}

	mockPDFService.EXPECT().
		ConvertToPDF(gomock.Any(), gomock.Eq(request)).
		Return(expectedResponse, nil)

	res, err := mockPDFService.ConvertToPDF(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, expectedResponse.PdfData, res.PdfData)

	mockPDFService.EXPECT().
		ConvertToPDF(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("some error"))

	_, err = mockPDFService.ConvertToPDF(context.Background(), &gen.ConvertToPDFRequest{})
	assert.Error(t, err)
	assert.Equal(t, "some error", err.Error())
}

func TestGetSavedPDF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPDFClient := mocks.NewMockPDFServiceClient(ctrl)
	apiClient := APIClient{pdfClient: mockPDFClient}

	req := &pdf.GetSavedPDFRequest{UserId: 1, Description: "test"}
	expectedResp := &pdf.GetSavedPDFResponse{PdfData: []byte("pdf data")}

	mockPDFClient.EXPECT().GetSavedPDF(gomock.Any(), gomock.Eq(req), gomock.Any()).Return(expectedResp, nil)

	resp, err := apiClient.GetSavedPDF(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockPDFClient.EXPECT().GetSavedPDF(gomock.Any(), gomock.Eq(req), gomock.Any()).Return(nil, errors.New("some error"))

	resp, err = apiClient.GetSavedPDF(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeletePDF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPDFClient := mocks.NewMockPDFServiceClient(ctrl)
	apiClient := APIClient{pdfClient: mockPDFClient}

	req := &pdf.DeletePDFRequest{UserId: 1, Description: "test"}
	expectedResp := &pdf.DeletePDFResponse{Message: "mess"}
	mockPDFClient.EXPECT().DeletePDF(gomock.Any(), gomock.Eq(req), gomock.Any()).Return(expectedResp, nil)

	resp, err := apiClient.DeletePDF(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockPDFClient.EXPECT().DeletePDF(gomock.Any(), gomock.Eq(req), gomock.Any()).Return(nil, errors.New("some error"))

	resp, err = apiClient.DeletePDF(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteAllPDF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPDFClient := mocks.NewMockPDFServiceClient(ctrl)
	apiClient := APIClient{pdfClient: mockPDFClient}

	req := &pdf.DeleteAllPDFRequest{UserId: 1}
	expectedResp := &pdf.DeleteAllPDFResponse{Message: "asd"}
	mockPDFClient.EXPECT().DeleteAllPDF(gomock.Any(), gomock.Eq(req), gomock.Any()).Return(expectedResp, nil)

	resp, err := apiClient.DeleteAllPDF(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockPDFClient.EXPECT().DeleteAllPDF(gomock.Any(), gomock.Eq(req), gomock.Any()).Return(nil, errors.New("some error"))

	resp, err = apiClient.DeleteAllPDF(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
