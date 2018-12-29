package main

import (
	jsonT "../jsonTools"
	gobTools "../gobTools"
	//"fmt"
	//"os"
	"fmt"
	"bytes"
)

var jsonTestData = map[string]interface{}{
	"one": "1",
	"two": 2,
	"ext": map[string]interface{}{"s_b": 1, "s_b_a": 2, "s_b_b": 3, "s_b_c": "qwe"},
	"is_exixt": true,
}


type testStruct struct {
	One 	string `json:"one"`
	Two		int	   `json:"two"`
	Ext		struct{
		S_b 	int		`json:"s_b"`
		S_b_a	int		`json:"s_b_a"`
		S_b_b	int		`json:"s_b_b"`
		S_b_c	string	`json:"s_b_c"`
	}
}

//func main() {
//
//	ret := new(interface{})
//	jsonT.JsonLoads(jsonTestData, ret)
//
//	retJson, _ := (*ret).(map[string]interface{})
//	fmt.Println("Json get one ret is", retJson["one"])
//	fmt.Println("Json get ext ret is", retJson["ext"])
//
//	//jsonT.JsonIterGetOneFromByte()
//	//
//	//jsonT.JsonIterMashal()
//	//
//	//jsonT.JsonIterGet()
//
//	jsonObj, err := jsonT.JsonIterLoad(jsonTestData)
//	if nil != err {
//		fmt.Println("JSON load fail err is", err)
//		os.Exit(1)
//	}
//
//	ret1 := jsonObj("ext", "s_b").ToInt()
//	fmt.Println("JSON get jsonTestData /ext/s_b/ ret is", ret1)
//}



func main() {
	var retInt		int
	var retBool		bool
	var retString 	string

	retInt = jsonTestData["two"].(int)
	retInt = jsonTestData["ext"].(map[string]interface{})["s_b"].(int)


	testMap := make(map[string]interface{})
	testMap["a"] = 1

	testMapSub := make(map[string]string)
	testMapSub["sub_a"] = "aaa"

	testMap["Ext"] = testMapSub

	retInt = testMap["a"].(int)

	retString = testMap["Ext"].(map[string]string)["sub_a"]

	fmt.Println(retInt, retBool, retString)

	jsonString, _ := jsonT.JsonDumps(jsonTestData)

	ts := testStruct{}
	jsonT.JsonLoads(jsonString, &ts)

	retString = ts.One
	fmt.Println( retString )

	gobEncode := gobTools.GobEncode(-1)
	gobRet, _ := gobEncode(ts)

	fmt.Println( "GOb encode ret is", gobRet.Len() )
	fmt.Println( "GOb encode string to buffer ret is", bytes.NewBufferString(gobRet.String()).Len() )
	fmt.Println( "GOb encode string ret is", gobRet.String() )


	gobDecode := gobTools.GObDecode()

	gobDecodeStruct := testStruct{}

	//err := gobDecode(gobRet, &gobDecodeStruct)
	err := gobDecode(*bytes.NewBufferString(gobRet.String()), &gobDecodeStruct)
	fmt.Println("Gob decode ret is", err, gobDecodeStruct.One)
}