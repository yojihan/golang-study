# Unbuffered vs Buffered Channels

## Reference

https://www.youtube.com/watch?v=uwEfn_yMplE

## Unbuffered Channel

- 채널은 기본적으로 unbuffered
- 다른 고루틴이 값을 받을 수 있는 상태가 될 떄까지 send operation이 전송을 block
- receive operation은 고루틴이 값을 받을 수 있는 상태가 될 때까지 데이터를 block

```go
make(chan int)
make(chan int, 0)
```

## Buffered Channel

- buffered channel은 특정 capacity(channel 생성시 지정)만큼 값을 일단 저장해둠
- send operation은 capacity가 full인 경우에만 block
- receive operation은 channel이 비어있는 경우에만 block

```go
make(chan int, 5)
```

## Example

### Unbuffered Channel 

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	orders := make(chan string) // Unbuffered channel

	// Customer placing orders
	go func() {
		for i := 1; i <= 5; i++ {
			order := fmt.Sprintf("Coffee order #%d", i)
			fmt.Println("🗒️ Placed:", order)
			orders <- order
		}
		close(orders)
	}()

	// Barista processing orders
	for order := range orders {
		fmt.Printf("🫘 Preparing: %s\n", order)
		time.Sleep(2 * time.Second)
		fmt.Printf("☕️ Served: %s\n", order)
	}
}
```

- result

```
🗒️ Placed: Coffee order #1
🗒️ Placed: Coffee order #2
🫘 Preparing: Coffee order #1
☕️ Served: Coffee order #1
🫘 Preparing: Coffee order #2
🗒️ Placed: Coffee order #3
☕️ Served: Coffee order #2
🫘 Preparing: Coffee order #3
🗒️ Placed: Coffee order #4
☕️ Served: Coffee order #3
🫘 Preparing: Coffee order #4
🗒️ Placed: Coffee order #5
☕️ Served: Coffee order #4
🫘 Preparing: Coffee order #5
☕️ Served: Coffee order #5
```

### Buffered Channel

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	orders := make(chan string, 3) // Buffered channel with a capacity of 3

	// Customers placing orders
	go func() {
		for i := 1; i <= 5; i++ {
			order := fmt.Sprintf("Coffee order #%d", i)
			fmt.Println("🗒️ Placed:", order)
			orders <- order
		}
		close(orders)
	}()

	// Barista processing orders
	for order := range orders {
		fmt.Printf("🫘 Preparing: %s\n", order)
		time.Sleep(2 * time.Second)
		fmt.Printf("☕️ Served: %s\n", order)
	}
}

```

- result

```
🗒️ Placed: Coffee order #1
🗒️ Placed: Coffee order #2
🗒️ Placed: Coffee order #3
🗒️ Placed: Coffee order #4
🗒️ Placed: Coffee order #5
🫘 Preparing: Coffee order #1
☕️ Served: Coffee order #1
🫘 Preparing: Coffee order #2
☕️ Served: Coffee order #2
🫘 Preparing: Coffee order #3
☕️ Served: Coffee order #3
🫘 Preparing: Coffee order #4
☕️ Served: Coffee order #4
🫘 Preparing: Coffee order #5
☕️ Served: Coffee order #5
```

## Why Buffered Channels?

- producer와 consumer를 디커플링
- 트래픽을 늘릴 수 있음
- 처리량 증가
- worker pool 개발시 유용
- 비동기 시그널링 또는 이벤트에 유용