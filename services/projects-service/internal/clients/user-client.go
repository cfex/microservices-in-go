package clients

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cfex/microservices-in-go/services/common/genproto"
	"github.com/cfex/microservices-in-go/services/projects-service/cmd/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client genproto.UserServiceClient
}

func NewUserClient(cfg *config.Config) (*UserClient, error) {
	conn, err := grpc.NewClient(cfg.Server.UsrGrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	return &UserClient{conn: conn, client: genproto.NewUserServiceClient(conn)}, nil
}

func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *UserClient) GetUserById(ctx context.Context, userId string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.client.GetUserById(ctx, &genproto.GetUserRequest{ID: userId})
	if err != nil {
		return "", fmt.Errorf("failed to get user by id: %w", err)
	}

	return res.ID, nil
}
