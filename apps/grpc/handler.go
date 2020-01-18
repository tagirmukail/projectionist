package grpc

import (
	"context"
	"projectionist/config"
	projProto "projectionist/proto"
	"projectionist/provider"
)

type ProjectionistServer struct {
	cfg         *config.Config
	dbProvider  provider.IDBProvider
	cfgProvider provider.IDBProvider
}

func NewProjectionistServer(cfgProvider, dbProvider provider.IDBProvider, cfg *config.Config) *ProjectionistServer {
	return &ProjectionistServer{
		cfg:         cfg,
		dbProvider:  dbProvider,
		cfgProvider: cfgProvider,
	}
}

func (p *ProjectionistServer) NewUser(ctx context.Context, r *projProto.UserRequest) (*projProto.UserResponse, error) {
	return &projProto.UserResponse{}, nil
}
