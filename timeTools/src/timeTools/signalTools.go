package goTools

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"reflect"
)

var ExitReflectFunc = make([]func()(), 0)
//var ExitReflectMap = make(map[string]interface{})
//var ExitNameMap map[string]string

func SetCatchSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	exitChan := make(chan int)

	go func() {
		select {
		case sgl := <-c:
			fmt.Println("Recive signal", sgl, "program will quit")
			exitChan <- 0
			return
		}
	} ()

	<-exitChan
	for _, value := range ExitReflectFunc {
		value()
	}

	os.Exit(1)
}


func SetExitDo(obj interface{}, describe string, values ...interface{}) {
	//ExitNameMap[funcName] = describe
	//ExitReflectMap[describe] = obj

	f := func() {
		objF := reflect.ValueOf(obj)
		if objF.Kind() != reflect.Func {
			//print( "not reflect func",  objF.Kind().String())
			return
		}

		args := []reflect.Value{}
		if nil != values {

			for _, value := range values {
				args = append(args, reflect.ValueOf(value))
			}
		} else {
			args = nil
		}

		fmt.Print(describe)
		objF.Call(args)
		//ret := objF.Call(args)
		//print( ret[0].String() )
	}

	ExitReflectFunc = append(ExitReflectFunc, f)
}