/**
 * @File Name: logTools.go
 * @Author: lizy
 * @Email:
 * @Create Date: 2018-09-26 14:30:09
 * @Last Modified: 2017-12-17 12:12:30
 * @Description: 文件操作
 */

package logTools

import (
	"os"
	"path"
	"fmt"
	"bufio"
	"strings"
	"runtime"
	"errors"
	"time"
	"io/ioutil"
)

var WriteFileBuffer = 4096


func SprintfLog(template string, values ...interface{}) string {
	_, file, line, ok := runtime.Caller(1)

	var logInfo string
	if ok {
		ret_array := strings.Split(file, "/")
		ret_array_len := len(ret_array)
		logInfo = fmt.Sprintf("File %s line %d %s\n", ret_array[ret_array_len -1], line, template)
		logInfo = fmt.Sprintf(logInfo, values...)
	} else {
		logInfo = fmt.Sprintf("File None line None get caller fail %s\n", template)
		logInfo = fmt.Sprintf(logInfo, values...)

	}

	return logInfo
	//fmt.Printf("Result is %d, %s, %d, %d", emm, ret_array[], line, ok)
	//
	//return ""
}


func FileExist(filePath, fileName string, createDir bool) (bool, error) {
	_, err := os.Stat(filePath)
	if nil != err {
		if !createDir {
			//fmt.Printf("No this directory %s, %s\n", filePath, err.Error())
			return false, err
		}

		if err = os.MkdirAll(filePath, 0755); nil != err {
			//fmt.Printf("Create this directory %s fail, %s\n", filePath, err.Error())
			return false, err
		}
	}

	abslolutePathFile := path.Join(filePath, fileName)
	if _, err = os.Stat(abslolutePathFile); nil != err {
		return false, err
	}

	return true, nil
}


func ReadFile(filePath, fileName string, trimSuffix bool) (bool, error, func (close bool) (bool, error, string)) {
	if ok, err := FileExist(filePath, fileName, false); !ok {
		return false, err, nil
	}

	absolutePathFile := path.Join(filePath, fileName)

	f, err := os.Open(absolutePathFile);
	if  nil != err {
		//fmt.Printf("Open file %s fail, %s\n", abslolutePathFile, err.Error() )
		return false, err, nil
	}

	rd := bufio.NewReader(f)

	return true, err, func (close bool) (bool, error, string) {
		// 当一行长度过长时，数据会被截断 isPrefix 为 true，不推荐使用此方法
		//byteLine, isPrefix, err := rd.ReadLine()

		if close {
			return false, nil, ""
		}

		line, err := rd.ReadString('\n')

		if nil != err && "EOF" == err.Error() {

			if "" == line {
				return false, err, line
			}

			return true, err, line
		}

		//if "" == line && nil != err {
		//	fmt.Printf("file %s is end %s\n", abslolutePathFile, err.Error() )
		//	return false, ""
		//}

		if nil != err {
			//fmt.Printf("Read file %s fail %s\n", abslolutePathFile, err.Error() )
			return false, err, line
		}

		if trimSuffix {
			return true, nil, strings.TrimSuffix(line, "\n")
		}

		return true, nil, line
	}
}


func ReadFileAll(filePath, fileName string) ([]byte, error) {
	absolutePathFile := path.Join(filePath, fileName)
	return ioutil.ReadFile(absolutePathFile)
}


// O_RDONLY int = syscall.O_RDONLY // open the file read-only.
// O_WRONLY int = syscall.O_WRONLY // open the file write-only.
// O_RDWR   int = syscall.O_RDWR   // open the file read-write.
// O_APPEND int = syscall.O_APPEND // append data to the file when writing.
// O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
// O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist
// O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
// O_TRUNC  int = syscall.O_TRUNC  // if possible, truncate file when opened.


func Rewrite(filePath, fileName string) (bool, error, func (content string, bClose, aClose bool) (int, error) ) {
	if ok, err := FileExist(filePath, fileName, true); !ok {
		return false, err, nil
	}

	abslolutePathFile := path.Join(filePath, fileName)

	f, err := os.OpenFile(abslolutePathFile,  os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if nil != err {
		return false, err, nil
	}

	writeHandle := bufio.NewWriterSize(f, WriteFileBuffer)

	return true, nil, func (content string, bClose, aClose bool) (int, error) {
		// 写数据之前 关闭
		if bClose {
			writeHandle.Flush()
			f.Close()
			return 0, nil
		}

		writeLen, err := writeHandle.WriteString(content)

		// 写数据之后 关闭
		if aClose {
			writeHandle.Flush()
			f.Close()
			return writeLen, err
		}

		return writeLen, err
	}
}


func Append(filePath, fileName string) (bool, error, func (content string, bClose, aClose bool) (int, error) ) {
	if ok, err := FileExist(filePath, fileName, true); !ok {
		return false, err, nil
	}

	abslolutePathFile := path.Join(filePath, fileName)

	f, err := os.OpenFile(abslolutePathFile,  os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if nil != err {
		return false, err, nil
	}

	writeHandle := bufio.NewWriterSize(f, WriteFileBuffer)

	return true, nil, func (content string, bClose, aClose bool) (int, error) {
		// 写数据之前 关闭
		if bClose {
			writeHandle.Flush()
			f.Close()
			return 0, nil
		}

		writeLen, err := writeHandle.WriteString(content)

		// 写数据之后 关闭
		if aClose {
			writeHandle.Flush()
			f.Close()
			return writeLen, err
		}

		return writeLen, err
	}
}


// support goroutine
// todo file close not realize
func MultipleAppend(basePath string, ) ( bool, error, func (file, content string) (error) ) {
	fileHandleMap := make( map[string]*os.File )
	fileIoHandleMap := make( map[string]*bufio.Writer )
	//fileIoHandleMap := make( map[string]func (content string, bClose, aClose bool) (int, error) )
	fileNameCh := make( chan string )
	fileHandleCh := make( chan *bufio.Writer )
	//fileHandleCh := make(chan func (content string, bClose, aClose bool) (int, error))

	FileHandleErrMap := make( map[string]error )

	go func () {
		for {
			select {
			case fileName := <-fileNameCh:
				value, ok := fileIoHandleMap[fileName]
				if ok {
					fileHandleCh <- value
				} else {
					f, err := os.OpenFile(fileName,  os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
					if nil != err {
						FileHandleErrMap[fileName] = err
						fileHandleCh <- nil
						continue
					}

					writeHandle := bufio.NewWriterSize(f, WriteFileBuffer)
					fileHandleMap[fileName] = f
					fileIoHandleMap[fileName] = writeHandle

					fileHandleCh <- writeHandle
				}
			default:
				;
			}


		}
	} ()

	// 创建目录
	if ok, err := FileExist(basePath, "test", true); !ok {
		return ok, err, nil
	}

	return true, nil, func (file, content string) (error) {
		if ok, err := FileExist(basePath, file, true); !ok {
			return err
		}

		absFile := path.Join( basePath, file )
		fileNameCh <- absFile
		fileHandle := <-fileHandleCh

		if nil == fileHandle {
			return FileHandleErrMap[absFile]
		}

		fileHandle.WriteString(content)
		return nil
	}
}


func GoChanTest() {
	testMap := map[string]string{"key0": "value0", "key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4"}
	keyChan := make(chan string)
	valueChan := make(chan string)

	go func() {
		for {
			select {
			case key := <-keyChan:
				//fmt.Println("recv key ", key)
				value, ok := testMap[key]
				if ok {
					valueChan <- value
				} else {
					valueChan <- ""
				}
			default:

			}
		}
	} ()

	fmt.Println("Run here 255")

	for i := 0; i < 3 ; i++  {
		fmt.Println("Run here 258")
		go func(numb int) {
			fmt.Println("recv i ", numb)
			for {
				key := fmt.Sprintf("key%d", numb)
				keyChan <- key
				select {
				case ret := <- valueChan:
					fmt.Printf("key %s value %s\n", key, ret)
					//time.Sleep(time.Second * 100)
				}
			}
		} ( i )
	}
}


// 测试协程中 Map 是否可读、可写
// 可读、不可写
func GoroutineMapReadT() {
	tMap := map[string]string{"qwe": "asd", "rty": "dfg"}

	for i := 0; i < 3; i++ {
		go func (i int) {
			for {
				// success
				//for key, value := range tMap {
				//	fmt.Println("range ", i, key, value)
				//}

				// fail error
				tMap["lalala"] = "biubiubiu"
			}
		} (i)
	}
}

// 2006-01-02 15:04:05
// 划分细度：H (小时)、D (天)
func TimeFileAppend(basePath, filePrefix, fileSuffix string, fileSplit byte, flushRate, closeDelay int) (error, func (string, bool) error) {
	// 创建目录
	if ok, err := FileExist(basePath, "", true); !ok {
		return err, nil
	}

	fileHandleMap := make( map[string]*os.File )
	fileIoHandleMap := make( map[string]*bufio.Writer )
	//fileIoHandleMap := make( map[string]func (content string, bClose, aClose bool) (int, error) )
	//fileNameCh := make( chan string )
	//fileHandleCh := make( chan *bufio.Writer )
	//fileHandleCh := make(chan func (content string, bClose, aClose bool) (int, error))

	//FileHandleErrMap := make( map[string]error )


	var sleepTime int
	var timeFmt string
	//var delayHour time.Duration
	if 'H' == fileSplit {
		sleepTime = 60 * 60
		timeFmt = "2006010215"
		//delayHour, _ = time.ParseDuration("1h")
	} else if 'D' == fileSplit {
		sleepTime = 60 * 60 * 24
		timeFmt = "20060102"
		//delayHour, _ = time.ParseDuration("24h")
	} else {
		errors.New("Para fileSplit canot find in H、D")
	}

	createFileErr := errors.New("")

	nowTimer := time.Now()
	nowMunite := nowTimer.Minute()
	nowSecond := nowTimer.Second()
	nowTimeStr := nowTimer.Format(timeFmt)

	fileName := filePrefix + nowTimeStr + fileSuffix
	absFileName := path.Join( basePath, fileName )

	//nextTimer := nowTimer.Add( delayHour )
	//nextTimeStr := nextTimer.Format(timeFmt)

	var first = true
	waitGoChan := make( chan int )
	go func() {
		for {

			var writeHandle *bufio.Writer
			f, err := os.OpenFile(absFileName,  os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			//fmt.Println("Open file ", absFileName, err)
			if nil != err {

				createFileErr = err
				goto SLEEP
				//FileHandleErrMap[fileName] = err
				//fileHandleCh <- nil
				//continue
			}

			writeHandle = bufio.NewWriterSize(f, 0) //WriteFileBuffer)
			fileHandleMap[nowTimeStr] = f
			fileIoHandleMap[nowTimeStr] = writeHandle

			if first {
				waitGoChan <- 1
				first = false
			}

			SLEEP:
			thisSleepT := sleepTime - nowMunite * 60 - nowSecond
			time.Sleep( time.Second * time.Duration(thisSleepT) )

			nowTimer = time.Now()
			nowMunite = nowTimer.Minute()
			nowSecond = nowTimer.Second()
			nowTimeStr = nowTimer.Format(timeFmt)
			fileName = filePrefix + nowTimeStr + fileSuffix
			absFileName = path.Join( basePath, fileName )
		}
	} ()

	// flush buffer to file
	var flushFirst = true
	flushSleepTime := flushRate - nowMunite % flushRate
	go func() {
		for {
			nowTimeStr := time.Now().Format(timeFmt)

			if flushFirst {
				//fmt.Println("flush sleep time is ", flushSleepTime)
				time.Sleep( time.Second * time.Duration(flushSleepTime) )
				flushFirst = false
			} else {
				//fmt.Println("flush sleep time is ", flushRate)
				time.Sleep( time.Second * time.Duration(flushRate) )
			}

			ioHandle, ok := fileIoHandleMap[nowTimeStr]
			if !ok {
				continue
			}

			ioHandle.Flush()
		}
	} ()

	// close file handle
	go func() {
		for {
			now := time.Now()

			nowTimeStr := now.Format(timeFmt)
			nowSecond := now.Minute() * 60 + now.Second()
			//nextHour := now.Add(delayHour)
			//ret := nextHour.Sub( now )

			closeSleepTime := sleepTime + closeDelay - nowSecond
			//fmt.Println("close sleep time is ", closeSleepTime)
			time.Sleep( time.Second * time.Duration(closeSleepTime) )

			fileHandle, ok := fileHandleMap[nowTimeStr]
			if !ok {
				continue
			}

			ioHandle, ok := fileIoHandleMap[nowTimeStr]
			if ok {
				ioHandle.Flush()
				delete( fileIoHandleMap, nowTimeStr )
			}

			fileHandle.Close()
			delete( fileHandleMap, nowTimeStr )
		}
	} ()

	// wait goroutine start
	//time.Sleep(time.Millisecond * 50)
	<-waitGoChan

	return nil, func(content string, close bool) error {
		if close {
			//fmt.Println("Defer enter there")
			for _, value := range fileIoHandleMap {
				value.Flush()
			}

			for _, value := range fileHandleMap {
				value.Close()
			}
		}

		//fmt.Println("write log file time", nowTimeStr)
		ioHandle, ok := fileIoHandleMap[nowTimeStr]
		if !ok {
			return createFileErr
		}

		_, err := ioHandle.WriteString( content  )
		return err
	}
}