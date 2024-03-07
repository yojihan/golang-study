# Mutexes and confinement

## Reference

https://www.youtube.com/watch?v=piutgCAvUjw

## Mutex Example

### without Mutex

```go
package main

import (
	"fmt"
	"sync"
)

func buyTicket(wg *sync.WaitGroup, userId int, remainingTickets *int) {
	defer wg.Done()

	if *remainingTickets > 0 {
		*remainingTickets-- // User purchases a ticket
		fmt.Printf("User %d purchased a ticket. Tickets remaining: %d\n", userId, *remainingTickets)
	} else {
		fmt.Printf("User %d found no ticket.\n", userId)
	}
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
```

```
User 7 purchased a ticket. Tickets remaining: 19
User 0 purchased a ticket. Tickets remaining: 11
User 13 purchased a ticket. Tickets remaining: 9
User 3 purchased a ticket. Tickets remaining: 18
User 11 purchased a ticket. Tickets remaining: 13
User 8 purchased a ticket. Tickets remaining: 16
User 4 purchased a ticket. Tickets remaining: 5
User 43 purchased a ticket. Tickets remaining: 0
User 15 found no ticket.
User 6 purchased a ticket. Tickets remaining: 10
User 16 found no ticket.
User 1 purchased a ticket. Tickets remaining: 7
User 46 found no ticket.
User 31 purchased a ticket. Tickets remaining: 4
User 9 purchased a ticket. Tickets remaining: 15
User 41 purchased a ticket. Tickets remaining: 2
User 42 purchased a ticket. Tickets remaining: 1
User 26 found no ticket.
User 14 purchased a ticket. Tickets remaining: 3
User 23 found no ticket.
User 10 purchased a ticket. Tickets remaining: 14
User 18 found no ticket.
User 32 found no ticket.
User 20 found no ticket.
User 22 found no ticket.
User 25 found no ticket.
User 38 found no ticket.
User 12 purchased a ticket. Tickets remaining: 12
User 47 found no ticket.
User 39 found no ticket.
User 48 found no ticket.
User 2 purchased a ticket. Tickets remaining: 6
User 35 found no ticket.
User 40 found no ticket.
User 45 found no ticket.
User 27 found no ticket.
User 44 found no ticket.
User 29 found no ticket.
User 30 found no ticket.
User 21 found no ticket.
User 33 found no ticket.
User 19 found no ticket.
User 49 purchased a ticket. Tickets remaining: 17
User 37 found no ticket.
User 24 found no ticket.
User 34 found no ticket.
User 36 found no ticket.
User 17 found no ticket.
User 5 purchased a ticket. Tickets remaining: 8
User 28 found no ticket.
```

### with Mutex

```go
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
	var tickets int = 50

	var wg sync.WaitGroup

	// Simulating a lot of users trying to buy tickets
	for userId := 0; userId < 200; userId++ {
		wg.Add(1)

		// But ticket for the user with ID userId
		go buyTicket(&wg, userId, &tickets)
	}

	wg.Wait()
}
```

```
User 0 purchased a ticket. Tickets remaining: 19
User 49 purchased a ticket. Tickets remaining: 18
User 7 purchased a ticket. Tickets remaining: 17
User 8 purchased a ticket. Tickets remaining: 16
User 9 purchased a ticket. Tickets remaining: 15
User 10 purchased a ticket. Tickets remaining: 14
User 11 purchased a ticket. Tickets remaining: 13
User 12 purchased a ticket. Tickets remaining: 12
User 13 purchased a ticket. Tickets remaining: 11
User 14 purchased a ticket. Tickets remaining: 10
User 15 purchased a ticket. Tickets remaining: 9
User 16 purchased a ticket. Tickets remaining: 8
User 17 purchased a ticket. Tickets remaining: 7
User 18 purchased a ticket. Tickets remaining: 6
User 19 purchased a ticket. Tickets remaining: 5
User 20 purchased a ticket. Tickets remaining: 4
User 21 purchased a ticket. Tickets remaining: 3
User 22 purchased a ticket. Tickets remaining: 2
User 23 purchased a ticket. Tickets remaining: 1
User 24 purchased a ticket. Tickets remaining: 0
User 25 found no ticket.
User 38 found no ticket.
User 39 found no ticket.
User 40 found no ticket.
User 41 found no ticket.
User 42 found no ticket.
User 43 found no ticket.
User 44 found no ticket.
User 45 found no ticket.
User 46 found no ticket.
User 47 found no ticket.
User 48 found no ticket.
User 3 found no ticket.
User 1 found no ticket.
User 2 found no ticket.
User 4 found no ticket.
User 31 found no ticket.
User 26 found no ticket.
User 27 found no ticket.
User 28 found no ticket.
User 29 found no ticket.
User 30 found no ticket.
User 34 found no ticket.
User 32 found no ticket.
User 33 found no ticket.
User 5 found no ticket.
User 35 found no ticket.
User 6 found no ticket.
User 36 found no ticket.
User 37 found no ticket.
```

## Confinement Example

```go
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
```

```
User 3 purchased a ticket. Tickets remaining: 19
User 1 purchased a ticket. Tickets remaining: 18
User 0 purchased a ticket. Tickets remaining: 17
User 24 purchased a ticket. Tickets remaining: 16
User 4 purchased a ticket. Tickets remaining: 15
User 5 purchased a ticket. Tickets remaining: 14
User 6 purchased a ticket. Tickets remaining: 13
User 7 purchased a ticket. Tickets remaining: 12
User 8 purchased a ticket. Tickets remaining: 11
User 9 purchased a ticket. Tickets remaining: 10
User 10 purchased a ticket. Tickets remaining: 9
User 11 purchased a ticket. Tickets remaining: 8
User 12 purchased a ticket. Tickets remaining: 7
User 13 purchased a ticket. Tickets remaining: 6
User 14 purchased a ticket. Tickets remaining: 5
User 15 purchased a ticket. Tickets remaining: 4
User 16 purchased a ticket. Tickets remaining: 3
User 17 purchased a ticket. Tickets remaining: 2
User 18 purchased a ticket. Tickets remaining: 1
User 19 purchased a ticket. Tickets remaining: 0
User 20 found no tickets.
User 21 found no tickets.
User 22 found no tickets.
User 23 found no tickets.
User 2 found no tickets.
User 36 found no tickets.
User 25 found no tickets.
User 26 found no tickets.
User 27 found no tickets.
User 28 found no tickets.
User 29 found no tickets.
User 30 found no tickets.
User 31 found no tickets.
User 32 found no tickets.
User 33 found no tickets.
User 34 found no tickets.
User 42 found no tickets.
User 37 found no tickets.
User 38 found no tickets.
User 39 found no tickets.
User 40 found no tickets.
User 41 found no tickets.
User 45 found no tickets.
User 43 found no tickets.
User 44 found no tickets.
User 47 found no tickets.
User 46 found no tickets.
User 48 found no tickets.
User 49 found no tickets.
User 35 found no tickets.
```