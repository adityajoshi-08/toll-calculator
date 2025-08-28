package main

import (
	"time"
	"toll-calculator/types"

	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LoggingMiddleware {
	return &LoggingMiddleware{
		next: next,
	}
}

func (l *LoggingMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {

		logrus.WithFields(logrus.Fields{
			"obu_id": data.OBUID,
			"lat":    data.Lat,
			"long":   data.Long,
			"took":   time.Since(start),
		}).Info("Producing to kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}
