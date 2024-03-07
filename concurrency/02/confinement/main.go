package main

import (
	"fmt"
	"sync"
)

func manageTicket(ticketChan chan int, doneChan chan bool, tickets *int) {
	for {
		select {
		case user := <-ticketChan:
			if *tickets > 0 {
				*tickets--
				fmt.Printf("User %d purchased a ticket. Tickets remaining: %d\n", user, *tickets)
			} else {
				fmt.Printf("User %d found no tickets.\n", user)
			}
		case <-doneChan:
			fmt.Printf("Tickets remaining: %d\n", *tickets)
		}

	}
}

func buyTicket(wg *sync.WaitGroup, ticketChan chan int, userId int) {
	defer wg.Done()

	ticketChan <- userId
}

func main() {
	var wg sync.WaitGroup
	tickets := 20
	ticketChan := make(chan int) // Channel for sending ticket purchase requests
	doneChan := make(chan bool)  // Channel for signaling the stop

	go manageTicket(ticketChan, doneChan, &tickets)

	for userId := 0; userId < 50; userId++ {
		wg.Add(1)
		go buyTicket(&wg, ticketChan, userId)
	}

	wg.Wait()
	doneChan <- true
}
