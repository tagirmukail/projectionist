package apps

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgraph-io/badger"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	projGrpc "projectionist/apps/grpc"
	"projectionist/config"
	"projectionist/consts"
	"projectionist/middleware"
	projPB "projectionist/proto"
	"projectionist/provider"
)

// RunGRPC - started grpc server
func RunGRPC(cfg *config.Config, sqlDB *sql.DB, badgerDB *badger.DB) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.GrpcPort))
	if err != nil {
		grpclog.Fatalf("listen error: %v", err)
	}

	cfgProvider, err := provider.NewCfgProvider(badgerDB)
	if err != nil {
		grpclog.Fatalf("cfg provider error: %v", err)
	}

	dbProvider := provider.NewDBProvider(sqlDB)

	grpc.EnableTracing = true
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				return nil
			}),
		grpc.UnaryInterceptor(middleware.ServerInterceptor(cfg)),
		grpc.ConnectionTimeout(3 * time.Second),
	}
	grpcServer := grpc.NewServer(opts...)
	defer grpcServer.GracefulStop()
	projPB.RegisterProjectionistServiceServer(
		grpcServer,
		projGrpc.NewProjectionistServer(
			cfgProvider,
			dbProvider,
			cfg,
		))
	reflection.Register(grpcServer)
	grpclog.Infof("start listen grpc server: %s:%d", cfg.Host, cfg.GrpcPort)
	err = grpcServer.Serve(listener)
	if err != nil {
		grpclog.Error(err)
	}
}

// RunGrpcApi - started HTTP reverse-proxy server
func RunGrpcApi(cfg *config.Config) {
	grpcServerEndpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.GrpcPort)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	marshalOption := runtime.WithMarshalerOption("*", &runtime.JSONPb{EmitDefaults: true})

	marshalOptionJson := runtime.WithMarshalerOption(consts.JsonOriginalType, &runtime.JSONPb{
		EmitDefaults: true,
		OrigName:     true,
	})

	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux(
		marshalOption,
		marshalOptionJson,
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := projPB.RegisterProjectionistServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		grpcServerEndpoint,
		opts,
	)
	if err != nil {
		grpclog.Errorf("projPB.RegisterProjectionistServiceHandlerFromEndpoint() error: %v", err)
		return
	}

	mux.Handle("/", gwmux)

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.GrpcApiPort))
	if err != nil {
		grpclog.Errorf("net.Listen error: %v", err)
		return
	}
	apiServer := http.Server{
		Addr: ln.Addr().String(),
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins(cfg.AccessAddresses),
		)(mux),
	}

	grpclog.Infof("start listen HTTP reverse-proxy server: %s:%d", cfg.Host, cfg.GrpcApiPort)
	err = apiServer.Serve(ln)
	if err != nil {
		grpclog.Errorf("apiServer.Serve error: %v", err)
		return
	}
}
