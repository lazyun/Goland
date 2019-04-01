package redisTools

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"reflect"
)

//const (
//	RedisNil 	= 0
//	RedisOK 	= 1
//	RedisErr 	= -1
//)

const (
	RedisRetFormat = "Value nil %t、Execute ret %t、Error exist %t、Key nil %t、effRet %d, %f、Error info %s"
)


type RedisCfg struct {
	RedisCluster []string `json:"redisCluster"`
}


type RedisRet struct {
	RedisNil		bool
	RedisOK 		bool
	RedisErr 		bool
	RedisExist		bool

	RedisRetInt		int64
	RedisRetFloat 	float64
	ErrInfo			error
}

type RedisTools struct {
	clusterClient	*redis.ClusterClient

	ClusterCfg		redis.ClusterOptions
	ClusterAddr		[]string

	RedisLogHandle  func(args ...interface{})
	RedisErr		error
}



func (this *RedisTools) SetConfig(redisAddr []string) error {
	this.ClusterAddr = redisAddr
	this.RedisLogHandle = func(args ...interface{}) {
		
	}
	return this.ClusterConnect()
}


func (this *RedisTools) ClusterConnect() error {
	this.ClusterCfg = redis.ClusterOptions{
		Addrs: this.ClusterAddr,
	}

	this.clusterClient = redis.NewClusterClient(&this.ClusterCfg)
	return this.clusterClient.Ping().Err()
}

func (this *RedisRet) RedisRetSpint() (s string) {
	s = fmt.Sprintf(RedisRetFormat, this.RedisNil,
		this.RedisOK, this.RedisErr, this.RedisExist,
			this.RedisRetInt, this.RedisRetFloat, this.ErrInfo)
	return
}


// 设置异常处理的接口
func (this *RedisTools) SetLogHandle(funcName interface{}) bool {
	refValue := reflect.ValueOf(funcName)
	if reflect.Func != refValue.Kind() {
		return false
	}

	f := func(args ...interface{}) {
		refArgs := []reflect.Value{}
		for _, value := range args {
			refArgs = append(refArgs, reflect.ValueOf(value))
		}

		refValue.Call(refArgs)
	}

	this.RedisLogHandle = f
	return true
}


func (this *RedisTools) SetLogHandleNew(funcName func (interface{})) {

	f := func(args ...interface{}) {
		content := ""
		if len(args) > 1 {
			content = fmt.Sprintf(args[0].(string), args[1:]...)
		}

		if len(args) == 1 {
			switch args[0].(type) {
			case error:
				content = args[0].(error).Error()
			default:
				content = fmt.Sprintf("%s", args)
			}
		}
		
		funcName(content)
	}

	this.RedisLogHandle = f
}



// 记录日志
func (this *RedisTools) SetLog(msg interface{}) {
	if nil == this.RedisLogHandle {
		return
	}

	this.RedisLogHandle(msg)
}

// ****************************************** 通用操作 ************************************************
// 判断一个 key 是否存在
// 返回结果为 true：redisRet.RedisErr：存在异常、redisRet.RedisOK 存在 or 不存在
func (this *RedisTools) Exist(key string) (redisRet RedisRet) {
	ret := this.clusterClient.Exists(key)
	
	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	if 1 == ret.Val() {
		redisRet.RedisOK = true
		return
	}

	return
}


// 设置一个 key 的过期时间 秒级别
// 返回结果为 true：RedisRet.RedisOK：获取成功、RedisRet.RedisErr：存在异常
func (this *RedisTools) ExpireKey(key string, expireTime time.Duration) (redisRet RedisRet) {
	ret := this.clusterClient.Expire(key, time.Second * expireTime)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo {
		redisRet.RedisOK = true
		return
	}

	this.SetLog(redisRet.ErrInfo)
	redisRet.RedisErr = true
	return
}


// 删除一个 key
// 返回结果为 true：RedisRet.RedisOK：删除成功、RedisRet.RedisErr：存在异常
func (this *RedisTools) Delete(key string) (redisRet RedisRet) {
	ret := this.clusterClient.Del(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true
	if 1 == ret.Val() {
		redisRet.RedisExist = true
		return
	}
	return
}


// 删除多个 key
// 返回结果为 true：RedisRet.RedisOK：获取成功、RedisRet.RedisErr：存在异常
func (this *RedisTools) DeleteKeys(key ...string) (redisRet RedisRet) {
	// todo coding...
	ret := this.clusterClient.Del(key ...)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		ret.Val(),
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo {
		redisRet.RedisOK = true
		return
	}

	this.SetLog(redisRet.ErrInfo)
	redisRet.RedisErr = true
	return
}


// 对 key 的 value +1
func (this *RedisTools) Incr(key string) (redisRet RedisRet) {
	ret := this.clusterClient.Incr(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		ret.Val(),
		0,
		ret.Err(),
		}

	if nil == redisRet.ErrInfo {
		redisRet.RedisOK = true
		return
	}

	this.SetLog(redisRet.ErrInfo)
	redisRet.RedisErr = true
	return
}


// 对 key 的 value 增加对应数值
func (this *RedisTools) IncrBy(key string, value int64) (redisRet RedisRet) {
	ret := this.clusterClient.IncrBy(key, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		ret.Val(),
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo {
		redisRet.RedisOK = true
		return
	}

	this.SetLog(redisRet.ErrInfo)
	redisRet.RedisErr = true
	return
}
// ****************************************** 字符串 ************************************************
// 获取一个字符串。
// value 支持 *string、*int64、*float64 格式
// 返回结果为 true：RedisRet.RedisNil：不存在 key-value、
// RedisRet.RedisOK：获取成功、RedisRet.RedisErr：存在异常
func (this *RedisTools) Get(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.Get(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if redis.Nil == redisRet.ErrInfo {
		redisRet.RedisNil = true
		return
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	switch value.(type) {
	case (*string):
		*value.(*string), redisRet.ErrInfo = ret.Result()

	case (*int64):
		*value.(*int64), redisRet.ErrInfo = ret.Int64()
	case (*float64):
		*value.(*float64), redisRet.ErrInfo = ret.Float64()
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true
	return
}


// 设置 key - value 并且设置 key 过期时间，0 代表不过期
// 返回结果为 true：RedisRet.RedisOK：设置成功、否则异常 RedisRet.RedisErr：存在异常
func (this *RedisTools) Set(key string, value interface{}, expire time.Duration) (redisRet RedisRet) {
	ret := this.clusterClient.Set(key, value, time.Second * expire)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo{
		redisRet.RedisOK = true
	} else {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
	}

	return
}


// 设置 key-value 当且仅当 key 不存在
// 返回结果为 true：RedisRet.RedisExist：存在 key-value 设置失败、
// RedisRet.RedisOK：获取成功、RedisRet.RedisErr：存在异常
func (this *RedisTools) SetNX(key string, value interface{}, expire time.Duration) (redisRet RedisRet) {
	ret := this.clusterClient.SetNX(key, value, time.Second *expire)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
	} else if !ret.Val() {
		redisRet.RedisExist = true
	} else {
		redisRet.RedisOK = true
	}

	return
}


// ****************************************** 链表 ************************************************
// 获取列表长度
//
func (this *RedisTools) LLen(key string) (int64, error) {
	ret := this.clusterClient.LLen(key)
	return ret.Result()
}

// 向列表左侧追加数据
// 返回结果为 true：RedisRet.RedisOK：设置成功、否则异常 RedisRet.RedisErr：存在异常
func (this *RedisTools) LPush(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.LPush(key, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo{
		redisRet.RedisOK = true
	} else {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
	}

	return
}


func (this *RedisTools) LPushKeys(key string, value ...interface{}) (redisRet RedisRet) {
	// todo coding
	ret := this.clusterClient.LPush(key, value ...)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo{
		redisRet.RedisOK = true
	} else {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
	}

	return
}


// 向列表右侧追加数据
// 返回结果为 true：RedisRet.RedisOK：设置成功、否则异常 RedisRet.RedisErr：存在异常
func (this *RedisTools) RPush(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.RPush(key, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo{
		redisRet.RedisOK = true
	} else {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
	}

	return
}


func (this *RedisTools) RPushKeys(key string, value ...interface{}) (redisRet RedisRet) {
	// todo coding
	ret := this.clusterClient.RPush(key, value ...)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil == redisRet.ErrInfo{
		redisRet.RedisOK = true
	} else {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
	}

	return
}


// 在列表中 找到值为 pivot 的位置，在其前方插入 value
// 返回结果：redisRet.RedisErr 存在异常、redisRet.RedisExist 查找的值是否存在、redisRet.RedisOK 操作是否成功
// 只有 redisRet.RedisExist、redisRet.RedisOK 同时为 true 代表插入成功
func (this *RedisTools) LInsertBefore(key string, pivot, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.LInsertBefore(key, pivot, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		ret.Val(),
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true

	// 当前列表长度
	if 1 <= ret.Val() {
		redisRet.RedisExist = true
		return
	}

	return
}


// 从左侧取出一个值并删除
// redisRet.RedisErr 是否存在异常、redisRet.RedisNil 是否存在 key、
// redisRet.RedisOK 是否取值成功
func (this *RedisTools) LPop(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.LPop(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if redis.Nil == redisRet.ErrInfo {
		redisRet.RedisNil = true
		return
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	switch value.(type) {
	case (*string):
		*value.(*string), redisRet.ErrInfo = ret.Result()

	case (*int64):
		*value.(*int64), redisRet.ErrInfo = ret.Int64()
	case (*float64):
		*value.(*float64), redisRet.ErrInfo = ret.Float64()
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true
	return
}


// 从左侧取出一个值并删除
// redisRet.RedisErr 是否存在异常、redisRet.RedisNil 是否存在 key、
// redisRet.RedisOK 是否取值成功
func (this *RedisTools) RPop(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.RPop(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if redis.Nil == redisRet.ErrInfo {
		redisRet.RedisNil = true
		return
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	switch value.(type) {
	case (*string):
		*value.(*string), redisRet.ErrInfo = ret.Result()

	case (*int64):
		*value.(*int64), redisRet.ErrInfo = ret.Int64()
	case (*float64):
		*value.(*float64), redisRet.ErrInfo = ret.Float64()
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true
	return
}


// ****************************************** 集合 ************************************************
// 向集合 Set 中添加一个 value
// 返回结果：redisRet.RedisErr 是否存在异常、
// redisRet.RedisOK 是否操作成功、redisRet.RedisNil 是否已经存在 true 不存在 false 存在
func (this *RedisTools) SAdd(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.SAdd(key, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		ret.Val(),
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true

	// 添加成功的数目
	if 1 <= ret.Val() {
		redisRet.RedisNil = true
		return
	}

	return
}


// 判断元素是否在集合中
// redisRet.RedisErr：是否有异常、redisRet.RedisOK 是否存在
func (this *RedisTools) SIsMembers(key string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.SIsMember(key, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	if ret.Val() {
		redisRet.RedisOK = true
	}

	return
}


// 以 Array[string] 的格式返回 Set 中的数据
func (this *RedisTools) SMembersArray(key string, value *[]string) (redisRet RedisRet) {
	ret := this.clusterClient.SMembers(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true

	*value = ret.Val()
	return
}


// 以 map[string]struct{} 的格式返回 Set 中的数据
func (this *RedisTools) SMembersMap(key string, value *map[string]struct{}) (redisRet RedisRet) {
	ret := this.clusterClient.SMembersMap(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true

	*value = ret.Val()
	return
}


// 判断两个集合的差集
// 返回结果：redisRet.RedisOK 是否操作成功
func (this *RedisTools) SExist(key1, key2 string, value *[]string) (redisRet RedisRet) {
	ret := this.clusterClient.SDiff(key1, key2)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true

	*value = ret.Val()
	return
}

// ****************************************** map ************************************************
// 在集合 key 中蛇者 field 和 value 键值对
// 返回结果：redisRet.RedisErr 是否存在异常、redisRet.RedisOK 是否操作成功
func (this *RedisTools) HSet(key, field string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.HSet(key, field, value)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	if ret.Val() {
		redisRet.RedisOK = true
	}

	return
}


// 获取集合中 field 对应 value 值
// value 支持 *string、*int64、*float64 格式
// 返回结果为 true：RedisRet.RedisNil：不存在 key-value、
// RedisRet.RedisOK：获取成功、RedisRet.RedisErr：存在异常
func (this *RedisTools) HGet(key, field string, value interface{}) (redisRet RedisRet) {
	ret := this.clusterClient.HGet(key, field)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if redis.Nil == redisRet.ErrInfo {
		redisRet.RedisNil = true
		return
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	switch value.(type) {
	case (*string):
		*value.(*string), redisRet.ErrInfo = ret.Result()

	case (*int64):
		*value.(*int64), redisRet.ErrInfo = ret.Int64()
	case (*float64):
		*value.(*float64), redisRet.ErrInfo = ret.Float64()
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true
	return
}


// 获取 key 下的所有键值对
func (this *RedisTools) HGetAll(key string, value *map[string]string) (redisRet RedisRet) {
	ret := this.clusterClient.HGetAll(key)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		0,
		0,
		ret.Err(),
	}

	if redis.Nil == redisRet.ErrInfo {
		redisRet.RedisNil = true
		return
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}


	*value = ret.Val()
	if 0 == len(*value) {
		redisRet.RedisNil = true
		return
	}

	redisRet.RedisOK = true
	return
}


// 对 hash 中 key 的值增加对应 incr 值
func (this *RedisTools) HIncrBy(key, field string, incr int64) (redisRet RedisRet) {
	ret := this.clusterClient.HIncrBy(key, field, incr)

	redisRet = RedisRet{
		false,
		false,
		false,
		false,
		ret.Val(),
		0,
		ret.Err(),
	}

	if nil != redisRet.ErrInfo {
		this.SetLog(redisRet.ErrInfo)
		redisRet.RedisErr = true
		return
	}

	redisRet.RedisOK = true
	return
}