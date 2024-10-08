package rpc

import (
	"context"
	"fmt"

	"github.com/infraboard/mcenter/clients/rpc"
	"github.com/infraboard/mcenter/clients/rpc/resolver"
	"github.com/infraboard/mflow/apps/task"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/infraboard/mcube/v2/grpc/middleware/exception"
	"github.com/infraboard/mcube/v2/ioc/config/log"
)

// NewClient todo
func NewClient() (*ClientSet, error) {
	log := log.Sub("sdk.mflow")

	mcenter := rpc.Config()
	ctx, cancel := context.WithTimeout(context.Background(), mcenter.Timeout())
	defer cancel()

	// 连接到服务
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s://%s?%s", resolver.Scheme, "mflow", mcenter.Resolver.ToQueryString()),
		grpc.WithPerRPCCredentials(mcenter.Credentials()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),

		// 将异常转化为 API Exception
		grpc.WithChainUnaryInterceptor(exception.NewUnaryClientInterceptor()),

		// Grpc Trace
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		return nil, err
	}

	return &ClientSet{
		conn: conn,
		log:  log,
	}, nil
}

// Client 客户端
type ClientSet struct {
	conn *grpc.ClientConn
	log  *zerolog.Logger
}

// 关闭GRPC连接
func (c *ClientSet) Stop() error {
	return c.conn.Close()
}

// Job Task 管理接口
func (s *ClientSet) JobTask() task.JobRPCClient {
	return task.NewJobRPCClient(s.conn)
}
