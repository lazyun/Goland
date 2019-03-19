package main

import (
	redisTools "../redisTools"
	"fmt"
)

//func main()  {
func mainsd()  {
	tools := new(redisTools.RedisTools)
	tools.SetConfig([]string{"127.0.0.1:6901", "127.0.0.1:6902"})
	tools.SetLogHandle(printErrInfo)

	redisKey := "a"
	var redisValue string
	//var redisValue int64
	ret := tools.Get(redisKey, &redisValue)
	if !ret.RedisOK {
		fmt.Println("Redis get value fail ", ret, redisValue)
	} else {
		fmt.Println("Redis get success value is", ret, redisValue)
	}

	// set
	ret = tools.Set("lalala", fmt.Sprint([]string{"biu", "biu", "~"}), 0)
	if ret.RedisOK {
		fmt.Println("Set key success!")
	} else {
		fmt.Println("Set key fail!", ret.ErrInfo)
	}

	// set not exist
	ret = tools.SetNX("lalala", fmt.Sprint([]string{"biu", "biu", "~"}), 0)
	if ret.RedisOK {
		fmt.Println("Set key success!", ret.ErrInfo)
	} else if ret.RedisErr {
		fmt.Println("Set key err!", ret.ErrInfo)
	} else {
		fmt.Println("Set key fail! key is early exist!", ret.ErrInfo)
	}

	// judge key exist
	ret = tools.Exist("lalala")
	fmt.Println("Key lalala is exist", ret.RedisOK, ret.ErrInfo)

	// delete a key
	ret = tools.Delete("lalala")
	fmt.Println("Key lalala is delete", ret.RedisOK, "key is exist", ret.RedisExist, ret.ErrInfo)

	// ********************************* array *********************************
	// list len
	listLen, err := tools.LLen("GoTestList")
	fmt.Println("key GoTestList len is", listLen, err )

	// delete a key
	ret = tools.Delete("GoTestList")
	fmt.Println("Key GoTestList is delete", ret.RedisOK, "key is exist", ret.RedisExist, ret.ErrInfo)

	// right push
	ret = tools.RPush("GoTestList", 2)
	fmt.Println("List GoTestList left push", ret.RedisOK, "error info", ret.RedisExist, ret.ErrInfo)

	// left push
	ret = tools.LPush("GoTestList", 5)
	fmt.Println("List GoTestList left push", ret.RedisOK, "error info", ret.RedisExist, ret.ErrInfo)

	// before insert
	ret = tools.LInsertBefore("GoTestList", 2, 11)
	fmt.Println("List GoTestList before insert", ret.RedisOK, "pivot exist", ret.RedisExist, "error info", ret.RedisExist, ret.ErrInfo)

	// before insert
	ret = tools.LInsertBefore("GoTestList", 12, 111)
	fmt.Println("List GoTestList before insert", ret.RedisOK, "pivot exist", ret.RedisExist, "error info", ret.RedisExist, ret.ErrInfo)

	// left pop
	var listValue int64
	ret = tools.LPop("GoTestList", &listValue)
	fmt.Println("LPop GoTestList value", listValue, ret.RedisRetSpint())

	ret = tools.LPop("GoTestList1", &listValue)
	fmt.Println("LPop GoTestList1 value", listValue, ret.RedisRetSpint())

	// ********************************* set *********************************
	// set add
	ret = tools.SAdd("GoTestSet", 123)
	fmt.Println("SAdd GoTestSet ret", ret.RedisRetSpint())

	// set membery array
	var valueArray []string
	ret = tools.SMembersArray("GoTestSet", &valueArray)
	fmt.Println("SMembersArray GoTestSet ret", valueArray, ret.RedisRetSpint())

	// set membery map
	var valueMap map[string]struct{}
	ret = tools.SMembersMap("GoTestSet", &valueMap)
	fmt.Println("SMembersArray GoTestSet ret", valueMap, valueMap["123"], ret.RedisRetSpint())

	// ********************************* map *********************************
	// hash map set
	ret = tools.HSet("GoTestMap", "key1", 1)
	fmt.Println("HSet GoTestMap ret", ret.RedisRetSpint())

	// hash map get
	var valueMapRet int64
	ret = tools.HGet("GoTestMap1", "key1", &valueMapRet)
	fmt.Println("HGet not exist key: GoTestMap1 key1 ret", valueMapRet, ret.RedisRetSpint())

	ret = tools.HGet("GoTestMap", "key", &valueMapRet)
	fmt.Println("HGet not exist key: GoTestMap key ret", valueMapRet, ret.RedisRetSpint())

	ret = tools.HGet("GoTestMap", "key1", &valueMapRet)
	fmt.Println("HGet GoTestMap key1 ret", valueMapRet, ret.RedisRetSpint())

	var valueMapRet1 string
	ret = tools.HGet("GoTestMap", "key1", &valueMapRet1)
	fmt.Println("HGet GoTestMap key1 ret", valueMapRet1, ret.RedisRetSpint())

	// hash map getall
	var valueMapGetAll map[string]string
	ret = tools.HGetAll("GoTestMap1", &valueMapGetAll)
	fmt.Println("HGetAll not exist key: GoTestMap1 ret", valueMapGetAll, ret.RedisRetSpint())

	ret = tools.HGetAll("GoTestMap", &valueMapGetAll)
	fmt.Println("HGetAll GoTestMap ret", valueMapGetAll, ret.RedisRetSpint())

	// hash incrby 10
	ret = tools.HIncrBy("GoTestMap", "key1", 10)
	fmt.Println("HIncrBy GoTestMap key1 ret", ret.RedisRetSpint())


}


func printErrInfo(msg string) {
	fmt.Println(msg)
}