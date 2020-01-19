package middleware

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"projectionist/config"
	"projectionist/consts"
	"projectionist/validate"
	"strings"
	"time"
)

func authorize(ctx context.Context, tokenSecretKey string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("retrieving metadata is failed")
	}

	authHeader, ok := md[strings.ToLower(consts.AuthorizationHeader)]
	if !ok {
		return fmt.Errorf("authorization token is not supplied")
	}

	if len(authHeader) == 0 {
		return fmt.Errorf("authorization token is empty")
	}

	err := validate.ValidateToken(authHeader[0], tokenSecretKey)
	if err != nil {
		return fmt.Errorf("authorize error: %v", err)
	}

	return nil
}

func ServerInterceptor(cfg *config.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (
		resp interface{}, err error) {
		start := time.Now()
		if info.FullMethod != "/projectionist.ProjectionistService/Login" {
			err := authorize(ctx, cfg.TokenSecretKey)
			if err != nil {
				return nil, err
			}
		}

		resp, err = handler(ctx, req)

		grpclog.Infof(
			"Request - Method:%s\tDuration:%v\tError:%v\n",
			info.FullMethod,
			time.Since(start),
			err,
		)

		return resp, err
	}
}
