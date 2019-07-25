package main

import (
	"fmt"
	"time"

	"../mgoTools"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	//var dasdas = map[string]int{}
	//fmt.Println("Get qwewqe value is", dasdas["qwewqe"])
	//return
	time.Sleep(1)

	mgoCfg := mgoTools.MgoCfg{
		MongoURL: "mongodb://lizy:%s@127.0.0.1:27017/admin",
		Password: "lizy123",
		WorkDB:   "tt",
	}

	mgoT := mgoTools.MgoTools{}
	if !mgoT.SetConfig(mgoCfg.MongoURL, mgoCfg.Password, mgoCfg.WorkDB) {
		fmt.Println("Monngo connect fail!\n error is", mgoT.MgoErr)
	}

	fmt.Println("MongoDB connect success!")

	var col = "rule611_t"

	queryConditions := bson.M{
		"sd_cid": 2,
		"value": bson.M{
			"$regex": "^asd",
		},
	}

	updateConditions := bson.M{
		"$set": bson.M{
			"a": "b",
		},
	}

	var ret = map[string]interface{}{}

	// err
	//var ret map[string]interface{}

	if nil == ret {
		fmt.Println("ret is null")
	} else {
		fmt.Println("ret is not null", len(ret))
	}

	mgoT.FindAndModify(col, queryConditions, updateConditions, ret, true, true)
	fmt.Println(ret)
	fmt.Println(mgoT.MgoErr)

	var allRet = []map[string]interface{}{}
	var findProjections = map[string]interface{}{
		"sd_cid": 0,
	}
	findRet := mgoT.Find(col, queryConditions, findProjections)
	err := findRet.All(&allRet)
	err1 := findRet.Close()
	fmt.Println(&allRet, err, err1)

	if nil == allRet {
		fmt.Println("allRet is null")
	} else {
		fmt.Println("allRet is not null", len(allRet))
	}

	objId := bson.NewObjectId()
	fmt.Println(objId)

	//fmt.Println( mgoT.Insert(col, bson.M{ "_id": objId, "sd_cid": 12, "value": "zxc#hehehe" }) )

	fmt.Println(mgoT.UpdateOne(col, bson.M{"sd_cid": 12, "value": "zxc#hehehe"}, bson.M{"$addToSet": bson.M{"asdad": objId}}))
	//fmt.Println( ttttMapValue() )

	var rett = map[string][]bson.ObjectId{}

	fmt.Println(mgoT.FindOne(col, bson.M{"sd_cid": 12, "value": "zxc#hehehe"}, rett, bson.M{"asdad": 1, "_id": 0}), rett)
}

func printLog(msg interface{}) {
	fmt.Println(msg)
}

func ttttMapValue() map[string][]string {
	var gValue = "123"
	var tMap = map[string][]string{
		"A": []string{"1"},
		"B": []string{"2"},
		"C": []string{"3"},
	}

	for key, value := range tMap {
		fmt.Println(key)
		if key == "3" {
			gValue = value[0]
			continue
		}

		tMap[key] = append(value, gValue)
	}

	return tMap
}
