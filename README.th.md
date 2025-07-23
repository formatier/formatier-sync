# ฟอร์มาเทียร์ ซิงค์: จัดการส่งข้อมูลล่วงหน้าในงานที่เป็น Asynchronous ได้ง่ายๆ

library นี้ คือตัวช่วยในการส่งข้อมูลข้าม goroutine โดยส่งมาเป็น **Future** คล้าย Promise ใน Javascript (js) หรือ Typescript (ts)

## การใช้งาน

### ใช้งานแบบง่าย

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

### ใช้งานกับ channel

```go
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
```

โดยปกติ ```future.Warp()``` จะทำงานแบบ async หมายความว่า function ทุกอย่างที่อยู่ใน future.Warp จพทำงานบนพื้นหลังทั้งหมด
เว้นแต่ คุณจะใช้ ```future.NewFuture()``` คุณจะได้รับ ```*future.Future[T]``` มา ซึ่งคุฯสามารถเรียกใช้ ```ft.SendValue()``` ได้

-----

> ## คำเตือน
> ห้ามใช้ ```ft.SendValue()``` ร่วมกับ ```future.Warp()``` เด็ดขาด

-----