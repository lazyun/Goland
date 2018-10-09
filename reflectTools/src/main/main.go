package main

import (
	reflectT "../reflectTools"
	"fmt"
	"time"
)

func main() {
	reflectT.ReflectTest(testFunc)

	reflectT.ReflectVarTypeValue("ASda")

	reflectT.ReflectFuncDo(testFunc)

	reflectT.ReflectFuncDo(printSm, "la~la~la~")

	chanCloseTest()
}


func chanCloseTest() {
	tChan := make(chan int)

	go func() {
		var sss int
		for {
			tChan <- sss
			sss += 1
			if 4 == sss {
				close(tChan)
				tChan = nil
				return
			}
			time.Sleep(time.Second * 1)
		}
	} ()

	go func() {
		for {
			select {
			case ret := <- tChan:
				fmt.Println("Chan receive value is", ret)
				time.Sleep(time.Second * 1)
			}
		}
	} ()

	time.Sleep(time.Second * 3600)
}


func testFunc() {
	fmt.Println("This is test Func skr~")
}


func printSm(sm interface{}) {
	fmt.Println("PrintSm Func print", sm)
}