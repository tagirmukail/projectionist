package controllers

import (
	"github.com/golang/mock/gomock"
	"projectionist/provider"
	"testing"
)

type Helper struct {
	provider     *provider.MockIDBProvider
	mockProvider *provider.MockIDBProviderMockRecorder
	ctrl         *gomock.Controller
}

func NewHelper(t *testing.T) *Helper {
	ctrl := gomock.NewController(t)
	mock := provider.NewMockIDBProvider(ctrl)

	return &Helper{
		provider:     mock,
		mockProvider: mock.EXPECT(),
		ctrl:         ctrl,
	}
}
