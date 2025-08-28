package main

import (
	"encoding/json"
	"time"
	"toll-calculator/aggregator/client"
	"toll-calculator/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   *client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, cl *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		panic(err)
	}
	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   cl,
	}, nil
}

func (kc *KafkaConsumer) Start() {
	kc.isRunning = true
	kc.readMessageLoop()
}

func (kc *KafkaConsumer) readMessageLoop() {
	for kc.isRunning {
		msg, err := kc.consumer.ReadMessage(time.Second)
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
				continue
			}
			logrus.Errorf("Kafka Consume error: %v (%v)\n", err, msg)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialisation error: %v", err)
			continue
		}
		distance, err := kc.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Distance calculation error: %v", err)
			continue
		}
		// _ = distance
		req := types.Distance{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			OBUID: data.OBUID,
		}
		if err := kc.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregator error: %v", err)
			continue
		}
		// logrus.Infof("Received message on %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
		// logrus.Infof("Calculated distance for OBU ID [%d]: %.2f km", data.OBUID, distance)
	}
}
