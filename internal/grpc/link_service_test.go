package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/0x0FACED/link-saver-telegram-bot/internal/grpc/mocks"
	"github.com/0x0FACED/proto-files/link_service/gen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockLinkServiceClient(ctrl)
	apiClient := APIClient{client: mockClient}

	req := &gen.GetLinkRequest{UrlId: 1}
	expectedResp := &gen.GetLinkResponse{GeneratedUrl: "https://example.com"}
	mockClient.EXPECT().GetLink(gomock.Any(), req).Return(expectedResp, nil)

	resp, err := apiClient.GetLink(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockClient.EXPECT().GetLink(gomock.Any(), req).Return(nil, errors.New("some error"))

	resp, err = apiClient.GetLink(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockLinkServiceClient(ctrl)
	apiClient := APIClient{client: mockClient}

	req := &gen.GetLinksRequest{UserId: 1}
	expectedResp := &gen.GetLinksResponse{Links: []*gen.Link{&gen.Link{GeneratedUrl: "https://example.com"}}}
	mockClient.EXPECT().GetLinks(gomock.Any(), req).Return(expectedResp, nil)

	resp, err := apiClient.GetLinks(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockClient.EXPECT().GetLinks(gomock.Any(), req).Return(nil, errors.New("some error"))

	resp, err = apiClient.GetLinks(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetAllLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockLinkServiceClient(ctrl)
	apiClient := APIClient{client: mockClient}

	req := &gen.GetAllLinksRequest{UserId: 1}
	expectedResp := &gen.GetAllLinksResponse{Links: []*gen.Link{&gen.Link{GeneratedUrl: "https://example.com"}}}
	mockClient.EXPECT().GetAllLinks(gomock.Any(), req).Return(expectedResp, nil)

	resp, err := apiClient.GetAllLinks(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockClient.EXPECT().GetAllLinks(gomock.Any(), req).Return(nil, errors.New("some error"))

	resp, err = apiClient.GetAllLinks(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestSaveLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockLinkServiceClient(ctrl)
	apiClient := APIClient{client: mockClient}

	req := &gen.SaveLinkRequest{UserId: 1, OriginalUrl: "https://example.com"}
	expectedResp := &gen.SaveLinkResponse{Success: true}
	mockClient.EXPECT().SaveLink(gomock.Any(), req).Return(expectedResp, nil)

	resp, err := apiClient.SaveLink(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockClient.EXPECT().SaveLink(gomock.Any(), req).Return(nil, errors.New("some error"))

	resp, err = apiClient.SaveLink(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestDeleteLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockLinkServiceClient(ctrl)
	apiClient := APIClient{client: mockClient}

	req := &gen.DeleteLinkRequest{LinkId: 1}
	expectedResp := &gen.DeleteLinkResponse{Success: true}
	mockClient.EXPECT().DeleteLink(gomock.Any(), req).Return(expectedResp, nil)

	resp, err := apiClient.DeleteLink(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)

	mockClient.EXPECT().DeleteLink(gomock.Any(), req).Return(nil, errors.New("some error"))

	resp, err = apiClient.DeleteLink(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
