package client

import (
	"context"
	"fmt"
	"toll-calculator/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint         string
	client types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fmt.Println("Connection established to ", endpoint)
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		Endpoint:         endpoint,
		client: c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}