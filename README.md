# Formatier Sync: Easily Manage Asynchronous Data Delivery with Futures

This library helps you send data across goroutines using a **Future** pattern, similar to Promises in JavaScript (JS) or TypeScript (TS).

## Usage

### Simple Usage

```go
func MessageFuture() *future.Future[string] {
	ft := future.Warp(
		func() (string, error) {
			return "Gorn", nil
		},
		nil,
	)

	return ft
}

func main() {
	msgFt := MessageFuture()

	msg, err := msgFt.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(msg)
}
```

### Usage with Channels

```go
package main

import (
	"fmt"
	"time"

	future "github.com/formatier/formatier-sync"
)

func MessageFuture() *future.Future[string] {
	msgChan := make(chan string)

	// Wait 10 seconds, then send a message to the (chan string) in the background.
	go func() {
		time.Sleep(time.Second * 10)
		msgChan <- "Hello World"
	}()

	ft := future.Warp(
		func() (string, error) {
			return <-msgChan, nil
		},
		nil,
	)

	return ft
}

func main() {
	msgFt := MessageFuture()

	msg, err := msgFt.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(msg)
}
```

Typically, `future.Warp()` operates **asynchronously**, meaning all functions within `future.Warp` will execute in the background.

However, if you use `future.NewFuture()`, you'll receive a `*future.Future[T]` which allows you to call `ft.SendValue()`.

> ## Warning
>
> Do **NOT** use `ft.SendValue()` with `future.Warp()`.
