package main

import (
	"context"
	"toll-calculator/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{svc: svc}
}

func (s *GRPCAggregatorServer) AggregateDistance(req *types.AggregateRequest) error {
	distance := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return s.svc.AggregateDistance(distance)
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
    distance := types.Distance{
        OBUID: int(req.OBUID),
        Value: req.Value,
        Unix:  req.Unix,
    }
    return &types.None{}, s.svc.AggregateDistance(distance)
}