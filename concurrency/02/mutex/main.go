package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex

func buyTicket(wg *sync.WaitGroup, userId int, remainingTickets *int) {
	defer wg.Done()
	mutex.Lock()

	if *remainingTickets > 0 {
		*remainingTickets-- // User purchases a ticket
		fmt.Printf("User %d purchased a ticket. Tickets remaining: %d\n", userId, *remainingTickets)
	} else {
		fmt.Printf("User %d found no ticket.\n", userId)
	}
	mutex.Unlock()
}

func main() {
	var tickets int = 20

	var wg sync.WaitGroup

	// Simulating a lot of users trying to buy tickets
	for userId := 0; userId < 50; userId++ {
		wg.Add(1)

		// But ticket for the user with ID userId
		go buyTicket(&wg, userId, &tickets)
	}

	wg.Wait()
}
