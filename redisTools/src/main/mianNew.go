package main

import (
	"fmt"
	RTN "../redisToolsNew"
	"os"
	"time"
)

func main() {
	fmt.Println("lalala")
	RTN.PrintFile()

	redisCfg := RTN.RedisCfg{
		[]string{ "127.0.0.1:6901"},
		false,
	}

	var ret RTN.RedisRet
	rc := RTN.RedisTools{}
	if ok, err := rc.Init(redisCfg); !ok {
		fmt.Println("Connect redis fail error is", err)
		os.Exit(1)
	}

	rc.Set("a", "b", time.Duration(600))

	var value = ""
	ret = rc.Get("lalala", &value)
	fmt.Printf("Get %+v\n", ret)

	ret = rc.GetSet("biubiu", "biubiubiu~", &value)
	fmt.Printf("GetSet %+v\n", ret)

	ret = rc.SetNX("a", "b", time.Duration(600))
	fmt.Printf("SetNX %+v\n", ret)

	var incrValue int64
	ret = rc.Incr("aaaa", &incrValue)
	fmt.Printf("Incr %+v %d\n", ret, incrValue)

	ret = rc.IncrBy("aaaaa", 1, &incrValue)
	fmt.Printf("IncrBy %+v %d\n", ret, incrValue)

	ret = rc.Exist("aaaaaqqqdd")
	fmt.Printf("Exist %+v\n", ret)

	ret = rc.Expire("aaaaaqqqdd", 100)
	fmt.Printf("Expire %+v\n", ret)

	ret = rc.Delete("aaaaaqqqdd")
	fmt.Printf("Delete %+v\n", ret)

	ret = rc.Delete("a")
	fmt.Printf("Delete %+v\n", ret)

	ret = rc.LLen("laaaalen", &incrValue)
	fmt.Printf("LLen %+v %d\n", ret, incrValue)

	ret = rc.LPush("laaaa", 123)
	fmt.Printf("LPush %+v\n", ret)

	ret = rc.LLen("laaaa", &incrValue)
	fmt.Printf("LLen %+v %d\n", ret, incrValue)

	ret = rc.LPop("lpopop", &value)
	fmt.Printf("LPop %+v %s\n", ret, value)


}