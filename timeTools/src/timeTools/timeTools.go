package timeTools

import (
	"time"
)


const DefFormat = "2006-01-02 15:04:05"


/*
	返回当前时间对象
*/
func NowTimer() time.Time {
	return time.Now()
}


/*
	返回当前时间 unix 时间戳
*/
func NowUnix() int64 {
	return time.Now().Unix()
}


/*
	返回当前时间格式化字符串
*/
func NowTimeS(formatS string) string {
	return time.Now().Format(formatS)
}


/*
	返回 unix 时间对应的时间对象
*/
func UnixToTimer(thisUnix int64) (time.Time) {
	return time.Unix(thisUnix, 0)
}


/*
	返回 unix 时间对应的时间对象
*/
func UnixToStr(thisUnix int64, foramtStr string) (string) {
	return time.Unix(thisUnix, 0).Format(foramtStr)
}


/*
	返回给定时间字符串、格式化字符串返回时间对象
*/
func SToTimer(timeStamp, formatStr string) (time.Time, error) {
	return time.Parse(formatStr, timeStamp)
}


/*
	返回给定时间字符串、格式化字符串返回 unix 时间
*/
func SToUnix(timeStamp, formatStr string) (int64, error) {
	var err			error
	var thisTimer	time.Time

	if thisTimer, err = time.ParseInLocation(formatStr, timeStamp, time.Local); err == nil{
		return thisTimer.Unix(), nil
	}

	return 0, err
}


func NowDelayNext(delay int) (int, error) {
	now := time.Now()
	delayHour, err := time.ParseDuration("1h")
	if nil != err {
		return 0, nil
	}

	nowSecond := now.Minute() * 60 + now.Second()
	nextHour := now.Add(delayHour)
	ret := nextHour.Sub( now )

	return int(ret.Seconds()) + delay - nowSecond, nil
}