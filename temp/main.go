package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"toll-calculator/aggregator/client"
	"toll-calculator/types"
)

func main() {
	c, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal("Error connecting xD:", err)
	}

	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		OBUID: 1,
		Value: 55.5,
		Unix: time.Now().Unix(),
	}); err != nil {
		log.Fatal("Error aggregating:", err)
	}

	fmt.Println("Aggregated successfully")
}