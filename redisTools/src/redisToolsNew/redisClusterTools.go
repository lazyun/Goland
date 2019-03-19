package redisToolsNew

import (
	"fmt"
	"github.com/go-redis/redis"
	"runtime"
	"strings"
	"time"
)

var (
	file = ""

)


func init() {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		return
	}

	file = f[strings.LastIndex(f, "/") + 1:]
}


func PrintFile() {
	fmt.Println(file)
}

type RedisTools struct {
	clusterClient 		*redis.ClusterClient
	RedisClusterCfg		redis.ClusterOptions

	ClusterCfg			[]string

	//logH				LogHandle
}


type RedisRet struct {
	// now command execute success
	ok		bool

	// value whether exist
	exist	bool

	// list len
	llen	int64

	//cmd 	interface{}
	err     error
}


type RedisCfg struct {
	Hosts		[]string	`json:"hosts"`
	ReadOnly	bool		`json:"read_only"`
}


//type redisLogHandle struct {
//	isPrint		bool
//}
//
//
//func (rlh *redisLogHandle) SetLog(msg string) {
//	if !rlh.isPrint {
//		return
//	}
//
//	fmt.Println(msg)
//}
//
//
//func (rlh *redisLogHandle) PrintLog(ok bool) {
//	rlh.isPrint = ok
//}
//
//
//type LogHandle interface {
//	SetLog(string)
//	PrintLog(bool)
//}


// connect to redis cluster
// para: redisCfg type RedisCfg
// retuen: connect success return (true, nil) else (false, error)
func (r *RedisTools) Init(redisCfg RedisCfg) (bool, error) {
	r.ClusterCfg = redisCfg.Hosts
	r.RedisClusterCfg = redis.ClusterOptions{
		Addrs: r.ClusterCfg,
		ReadOnly: redisCfg.ReadOnly,
	}

	//r.logH = &redisLogHandle{ false }

	r.clusterClient = redis.NewClusterClient(&r.RedisClusterCfg)
	if _, err := r.clusterClient.Ping().Result(); nil != err {
		return false, err
	}

	return true, nil
}


// set redis err log redirct
//func (r *RedisTools) SetLoghandle(l LogHandle) {
//	r.logH = l
//}


// print log
//func (r *RedisTools) PrintInfo(ok bool) {
//	r.logH.PrintLog(ok)
//}


/******************************************************************** Currency Op *********************************************************************************/
// ok true cmd execute success
// exist true or false
func (r *RedisTools) Exist(key string) RedisRet {
	cmd := r.clusterClient.Exists(key)
	ret, err := cmd.Result()

	//fmt.Println(ret, err)
	if nil != err {
		return RedisRet{ ok:false }
	}

	if 0 == ret {
		return RedisRet{ ok:true, exist:false }
	}

	return RedisRet{ ok:true, exist:true }
}


func (r *RedisTools) Expire(key string, expire time.Duration) RedisRet {
	cmd := r.clusterClient.Expire(key, expire * time.Second)
	ret, err := cmd.Result()

	//fmt.Println(ret, err)
	if nil != err {
		return RedisRet{ ok:false }
	}

	if ret {
		return RedisRet{ ok:true, exist:true }
	}

	return RedisRet{ ok:false, exist:false }
}


func (r *RedisTools) Delete(key string) RedisRet {
	cmd := r.clusterClient.Del(key)
	ret, err := cmd.Result()

	//fmt.Println(ret, err)
	if nil != err {
		return RedisRet{ ok:false }
	}

	if 0 == ret {
		return RedisRet{ ok:true, exist:false }
	}

	return RedisRet{ ok:true, exist:true }
}


/******************************************************************** String Op *********************************************************************************/
// redis set key value expireTime
// return RedisRet.ok == true success else false err = RedisRet.err
func (r *RedisTools) Set(key string, value interface{}, expireTime time.Duration) RedisRet {
	cmd := r.clusterClient.Set(key, value, expireTime * time.Second)
	if _, err := cmd.Result(); nil != err {
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true }
}


// ok && exist all true mean get success
// exist false mean key dont exist
// ok true mean redis run cmd success
// ok false mean redus run cmd false
func (r *RedisTools) Get(key string, value *string) RedisRet {
	var ret string
	var err error
	cmd := r.clusterClient.Get(key)
	if ret, err = cmd.Result(); nil != err {
		if err == redis.Nil {
			return RedisRet{ ok:true, exist:false }
		}
		return RedisRet{ ok:false, err:err }
	}

	*value = ret
	return RedisRet{ ok:true, exist:true }
}


// same Get
func (r *RedisTools) GetSet(key string, value interface{}, ret *string) RedisRet {
	var err error
	cmd := r.clusterClient.GetSet(key, value)
	if *ret, err = cmd.Result(); nil != err {
		if err == redis.Nil {
			return RedisRet{ ok:true, exist:false }
		}
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true, exist:true }
}


// same Get
func (r *RedisTools) SetNX(key string, value interface{}, expireTime time.Duration) RedisRet {
	cmd := r.clusterClient.SetNX(key, value, expireTime * time.Second)
	ok, err := cmd.Result()

	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	if ok {
		return RedisRet{ ok:true, exist:false }
	}

	return RedisRet{ ok:true, exist:true }
}


// op success ok true else false
func (r *RedisTools) Incr(key string, ret *int64) RedisRet {
	var err error
	cmd := r.clusterClient.Incr(key)
	*ret, err = cmd.Result()

	//fmt.Println(*ret, err)

	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true }
}


func (r *RedisTools) IncrBy(key string, value int64, ret *int64) RedisRet {
	var err error
	cmd := r.clusterClient.IncrBy(key, value)
	*ret, err = cmd.Result()

	//fmt.Println(*ret, err)

	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true }
}


func (r *RedisTools) IncrByFloat(key string, value float64, ret *float64) RedisRet {
	var err error
	cmd := r.clusterClient.IncrByFloat(key, value)
	*ret, err = cmd.Result()

	//fmt.Println(*ret, err)

	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true }
}


/******************************************************************** List Op *********************************************************************************/
func (r *RedisTools) LLen(key string, value *int64) RedisRet {
	cmd := r.clusterClient.LLen(key)
	ret, err := cmd.Result()

	//fmt.Println(ret, err)
	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	*value = ret

	if 0 == ret {
		return RedisRet{ ok:true, exist:false }
	}

	return RedisRet{ ok:true, exist:true }
}


func (r *RedisTools) LPush(key string, value interface{}) RedisRet {
	cmd := r.clusterClient.LPush(key, value)
	ret, err := cmd.Result()

	//fmt.Println(ret, err)
	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true, llen:ret }
}


func (r *RedisTools) RPush(key string, value interface{}) RedisRet {
	cmd := r.clusterClient.RPush(key, value)
	ret, err := cmd.Result()

	//fmt.Println(ret, err)
	if nil != err {
		return RedisRet{ ok:false, err:err }
	}

	return RedisRet{ ok:true, llen:ret }
}


func (r *RedisTools) LPop(key string, value *string) RedisRet {
	cmd := r.clusterClient.LPop(key)
	ret, err := cmd.Result()

	if nil == err {
		*value = ret
		return RedisRet{ ok:true, exist:true }
	}

	if redis.Nil == err {
		return RedisRet{ ok:true, exist:false }
	}

	return RedisRet{ ok:false, err:err }
}


func (r *RedisTools) RPop(key string, value *string) RedisRet {
	cmd := r.clusterClient.RPop(key)
	ret, err := cmd.Result()

	if nil == err {
		*value = ret
		return RedisRet{ ok:true, exist:true }
	}

	if redis.Nil == err {
		return RedisRet{ ok:true, exist:false }
	}

	return RedisRet{ ok:false, err:err }
}