package main

import (
	"fmt"
	"time"

	future "github.com/formatier/formatier-sync"
)

func MessageFuture() *future.Future[string] {
	msgChan := make(chan string)

	// รอ 10 วินาที แล้วส่งข้อความเข้า (chan string) โดยทำงานในพื้นหลัง
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
