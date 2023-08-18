package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var receiveOrderCh = make(chan order)
	var validOrderCh = make(chan order)
	var invalidOrderCh = make(chan invalidOrder)

	go receiveOrder(receiveOrderCh)
	go validateOrder(receiveOrderCh, validOrderCh, invalidOrderCh)

	wg.Add(1)

	go func(validOrderCh <-chan order, invinvalidOrderCh <-chan invalidOrder) {
	loop:
		for {
			select {
			case order, ok := <-validOrderCh:
				if ok {
					fmt.Printf("Valid order received: %v\n", order)
				} else {
					break loop
				}

			case order, ok := <-invalidOrderCh:
				if ok {
					fmt.Printf("Invalid order was received: %v. Issue: %v\n", order.Order, order.err)
				} else {
					break loop
				}
			}
		}
		wg.Done()
	}(validOrderCh, invalidOrderCh)
	wg.Wait()

}

func validateOrder(in <-chan order, out chan<- order, errCh chan<- invalidOrder) {

	for order := range in {
		if order.Quantity <= 0 {
			errCh <- invalidOrder{Order: order, err: errors.New("quantity must be greater than 0")}
		} else {
			out <- order
		}
	}
	close(out)
	close(errCh)
}

func receiveOrder(out chan<- order) {

	for _, rawOrder := range rawOrders {
		var newOrder order
		err := json.Unmarshal([]byte(rawOrder), &newOrder)
		if err != nil {
			log.Print(err)
			continue
		}
		out <- newOrder
	}
	close(out)
}

var rawOrders = []string{
	"{\"productCode\": 1111, \"quantity\": 5, \"status\": 1}",
	"{\"productCode\": 2222, \"quantity\": 42.3, \"status\": 1}",
	"{\"productCode\": 3333, \"quantity\": 35, \"status\": 1}",
	"{\"productCode\": 4444, \"quantity\": 8, \"status\": 1}",
}
