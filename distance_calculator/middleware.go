package main

import (
	"time"
	"toll-calculator/types"

	"github.com/sirupsen/logrus"
)

type LogMiddleWare struct {
	next CalculatorServicer
}

func NewLogMiddleWare(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleWare{next: next}
}

func (lm *LogMiddleWare) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"distance": dist,
		}).Info("calculate distance")
	}(time.Now())
	dist, err = lm.next.CalculateDistance(data)
	return
}
