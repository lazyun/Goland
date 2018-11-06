package main

import (
	"fmt"
	MgoT "../mgoTools"
	"os"
	"gopkg.in/mgo.v2/bson"
)


func main() {
	fmt.Println("some")

	col := "go_test2"
	mgoT := new(MgoT.MgoTools)
	mgoUrl := "mongodb://wc_applet_rw:%s@s-2ze6c71dc69d0fc4.mongodb.rds.aliyuncs.com:3717,s-2zed48c86d79ee04.mongodb.rds.aliyuncs.com:3717/wc_applet"
	mgoPass := "DUw@87weh43gy#%WC"
	mgoDB := "wc_applet"
	ok := mgoT.SetConfig(mgoUrl, mgoPass, mgoDB)
	if !ok {
		fmt.Println(mgoT.MgoUrlStr)
		fmt.Printf("MongoDB connect fail err is %s \n", mgoT.MgoErr)
		os.Exit(1)
	}

	fmt.Println("MongoDB connect success!")

	// set log handle
	mgoT.SetLogHandle(PrintLog)

	// judge col is exist
	ok = mgoT.ColIsExist(col)
	if ok {
		fmt.Println("The collection ", col, " exist!")
	} else if nil != mgoT.MgoErr {
		fmt.Println("Use mongodb occur err ", mgoT.MgoErr)
	} else {
		fmt.Println("The collection ", col, "not exist!")
		//ok = mgoT.CreateCol(col)
		//
		//if !ok {
		//	fmt.Println("The collection ", col, "create fail!")
		//	panic("emm~ exit!")
		//}
	}

	// 创建索引
	//ok = mgoT.CreateIndexNormal(col, "start", true, true, false)
	//if ok {
	//	fmt.Println("集合 go_test 创建普通索引 start 成功！")
	//}
	//
	//ok = mgoT.CreateIndexHashed(col, "end", false)
	//if ok {
	//	fmt.Println("集合 go_test 创建哈希索引 end 成功！")
	//}
	//
	//var ret map[string]interface{}
	//ok = mgoT.ShardingCol("wc_applet", col, "end", true, &ret)
	//fmt.Println( "Sharding collection ret is ", ok, ret )

	//insertValues := bson.M{
	//	"start": "biu~biu~biu~",
	//	"end": "ying~ying~ying~",
	//	"q": 1,
	//	"w": 2,
	//	"e": 3,
	//	"r": 4,
	//	"t": 5,
	//}

	//insertValues1 := bson.M{}
	//
	//insertValues1["start"] = "biu~biu~biu~1"
	//insertValues1["end"] = "ying~ying~ying~1"
	//insertValues1["q"] = 1
	//insertValues1["w"] = 2
	//insertValues1["e"] = 3
	//insertValues1["r"] = 4
	//insertValues1["t"] = 5
	//insertValues1["y"] = 6
	//insertValues1["u"] = 7

	insertValues2 := bson.D{
		{"start", "biu~biu~biu~1"},
		{"end", "ying~ying~ying~1"},

		{"q", 1},
		{"w", 2},
		{"e", 3},
		{"r", 4},
		{"t", 5},
		{"y", 6},

		{
			"ext_info",
			bson.D{
				{"keySub1", "valueSub1"},
				{"keySub2", "valueSub2"},
				{"keySub3", "valueSub3"},
				{"keySub4", "valueSub4"},
				{"keySub5", "valueSub5"},

			},
		},
	}



	ok = mgoT.Insert(col, insertValues2)
	if ok {
		fmt.Println("insert success", insertValues2)
	} else if mgoT.MgoInsertDup {
		fmt.Println("insert data duplicate", insertValues2, mgoT.MgoErr)
	} else {
		fmt.Println("insert fail", insertValues2, mgoT.MgoErr)
	}

	// find
	query := bson.M{
		"end": "ying~ying~ying~1",
	}

	query_projections := bson.M{
		"start": 1,
		"_id": 0,
	}
	var deviceInfoOne map[string]interface{}
	ok = mgoT.FindOne("go_test", query, &deviceInfoOne, query_projections)
	fmt.Println( deviceInfoOne, ok, mgoT.MgoErr )

	// update
	//updateConditions := bson.M{
	//	"$set": bson.M{
	//		"start": "biu~biu~biu~2",
	//	},
	//}
	//ok = mgoT.UpdateOne(col, query, updateConditions)
	//if ok {
	//	fmt.Println("Update success", mgoT.MgoUpdateInfo)
	//}

	// delete
	//ok = mgoT.DeleteOne(col ,query)
	//if ok {
	//	fmt.Println("Delete one success!")
	//}
}


func userHandleTest(handle interface{}) (func(string, interface{}) (bool)) {
	afaf, ok := handle.(func(string, interface{}) (bool) )
	if ok {
		return afaf
	}

	fmt.Println(ok)
	return nil
}


func PrintLog(err error) {
	fmt.Println( "Reflect print err is ", err )
}