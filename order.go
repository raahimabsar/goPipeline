package main

import "fmt"

type order struct {
	ProductCode int
	Quantity    float64
	Status      orderStatus
}

type invalidOrder struct {
	Order order
	err   error
}

type orderStatus int

const (
	none orderStatus = iota
	new
	received
	reserved
	filled
)

func (o order) String() string {
	return fmt.Sprintf("Product code: %v, Quantity: %v, Status: %v", o.ProductCode, o.Quantity, orderStatusToText(o.Status))
}

func orderStatusToText(o orderStatus) string {
	switch o {
	case none:
		return "none"
	case new:
		return "new"
	case received:
		return "received"
	case filled:
		return "filled"
	}
	return "invalid status"
}
