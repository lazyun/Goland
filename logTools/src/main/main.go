package main

import (
	logT "../logTools"
	"os"
	"fmt"
	//"path"
	"time"
	"os/signal"
	"syscall"
)


func main() {
	//logT.GoChanTest()
	//logT.GoroutineMapReadT()

	//time.Sleep(time.Second * 30)
	//fmt.Println("test path.join /qwe/ert/, asdf/qwer/koior.log", path.Join("/qwe/ert/", "asdf/qwer/koior.log"))
	//os.Exit(1)

	// 读取文件
	ok, err, fileHandle := logT.ReadFile("/Users/phoenix/Documents/myGithub/Goland/logTools/docs", "test.log", true)
	if !ok {
		fmt.Printf("File read fail error is %s", err)
		os.Exit(1)
	}

	for {
		ok, err, line := fileHandle(false)
		if !ok {
			fileHandle(true)
			fmt.Printf("File read end is %s\n", err)
			break
		}

		fmt.Printf("Read line is %s\n\n", line)
	}

	// 重写文件
	ok, err, rewriteHandle := logT.Rewrite("/Users/phoenix/Documents/myGithub/Goland/logTools/docs", "test1.log")
	if !ok {
		fmt.Printf("File write fail error is %s", err)
		os.Exit(1)
	}

	//for i := 0; i < 5; i++ {
	//
	//	content := fmt.Sprintf("This is %d\n", i)
	//	writeHandle(content, false, false)
	//
	//	if 4 == i {
	//		writeHandle("", true, false)
	//	}
	//}

	for i := 0; i < 5; i++ {

		content := fmt.Sprintf("This is %d\n", i)
		rewriteHandle(content, false, false)
	};  rewriteHandle("", false, true)

	// 追加写文件
	ok, err, appendHandle := logT.Append("/Users/phoenix/Documents/myGithub/Goland/logTools/docs", "test1.log")
	if !ok {
		fmt.Printf("File write fail error is %s", err)
		os.Exit(1)
	}

	//for i := 0; i < 5; i++ {
	//
	//	content := fmt.Sprintf("This is %d\n", i)
	//	writeHandle(content, false, false)
	//
	//	if 4 == i {
	//		writeHandle("", true, false)
	//	}
	//}

	for i := 0; i < 5; i++ {

		content := fmt.Sprintf("This is append %d\n", i)
		appendHandle(content, false, false)
	};  appendHandle("", false, true)

	content := logT.SprintfLog("Run finished %s", "la~la~la~\n")
	fmt.Println(content)


	// 按时间区分文件
	timeFile, err := logT.TimeFileAppend("/Users/phoenix/Documents/myGithub/Goland/logTools/logs/", "li.", ".log", 'H', 600, 900)
	if nil != err {
		fmt.Println("Create time append file fail ", err)
		os.Exit(1)
	}

	defer timeFile("", true)

	var temporary int
	var exitChan = make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	okExitChan := make(chan int)

	go func() {
		select {
		case sgl := <-exitChan:
			fmt.Println("Recive signal", sgl, "program will quit")
			okExitChan <- 1
			return
		}
	} ()

	FOR:
	for {
		content := fmt.Sprintf("This is %d data\n", temporary)

		err = timeFile(content, false)

		fmt.Println( "Time app write err ", err )

		temporary += 1

		//time.Sleep(time.Millisecond * 50)

		select {
		case <-okExitChan:
			fmt.Println( "selecr exit signal\n " )
			break FOR
		default:
			time.Sleep(time.Second * 60)
			continue
		}
	}

	fmt.Println( "The program over\n " )
}
