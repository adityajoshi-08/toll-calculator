package main

import (
	"fmt"
	"log"
	"toll-calculator/aggregator/client"
)

const kafkaTopic = "obu_data"

func main() {
	var (
		err error
		svc CalculatorServicer
	)

	// httpClient := client.NewHTTPClient("http://localhost:3000/aggregate")
	grpcClient, err := client.NewGRPCClient("http://localhost:3000/aggregate")
	if err != nil {
		log.Fatal(err)
	}

	svc = NewCalculatorService()
	svc = NewLogMiddleWare(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	fmt.Println("Starting Kafka Consumer...")
}
