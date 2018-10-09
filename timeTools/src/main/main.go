package main

import (
	timeT "../timeTools"
	"fmt"
	"time"
)


func main() {
	ret, err := timeT.NowDelayNext(600)
	fmt.Println("ret is ", ret, err)

	//exitSignal := timeT.SetCatchSignal()
	//<-exitSignal

	//timeT.SetExitDo(timeT.NowTimeS, "lalala", timeT.DefFormat)

	go func() {
		for {
			fmt.Println("test is exit")
			time.Sleep(time.Second * 1)
		}

	} ()

	timeT.SetExitDo(printNowTimeStr, "biu~biu~biu~\t")
	timeT.SetCatchSignal()
}


func printNowTimeStr() {
	fmt.Println("now time is ", timeT.NowTimeS(timeT.DefFormat))
}