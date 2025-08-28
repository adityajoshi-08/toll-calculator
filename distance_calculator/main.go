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

	svc = NewCalculatorService()
	svc = NewLogMiddleWare(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient("http://localhost:3000/aggregate"))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	fmt.Println("Starting Kafka Consumer...")
}
