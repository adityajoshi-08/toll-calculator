package main

import (
	"fmt"
	"toll-calculator/types"
)

const basePricePerKm = 0.5

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (ia *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("Processing and insertign distance in the storage", distance)
	return ia.store.Insert(distance)
}

func (ia *InvoiceAggregator) CalculateInvoice(obuID int) (types.Invoice, error) {
	dist, err := ia.store.Get(obuID)
	if err != nil {
		return types.Invoice{}, fmt.Errorf("error getting distance for OBU ID %d: %v", obuID, err)
	}
	invoice := &types.Invoice{
		OBUID: obuID,
		TotalDistance: dist,
		TotalAmount: dist * basePricePerKm,
	}
	return *invoice, nil
}
