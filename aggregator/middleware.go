package main

import (
	"time"
	"toll-calculator/types"

	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{next: next}
}

func (lm *LogMiddleware) CalculateInvoice(obuID int) (invoice types.Invoice, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"obuID": obuID,
			"totalDistance": invoice.TotalDistance,
			"totalAmount": invoice.TotalAmount,
		}).Info("Calculate Invoice")
	}(time.Now())

	invoice, err = lm.next.CalculateInvoice(obuID)
	return
}

func (lm *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info()
	}(time.Now())

	err = lm.next.AggregateDistance(distance)
	return
}
