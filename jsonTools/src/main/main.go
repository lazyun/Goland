package main

import (
	jsonT "../jsonTools"
	"fmt"
	"os"
)

const jsonTestData = `{"one": "1", "two": 2, "ext": {"s_b": 1, "s_b_a": 2, "s_b_b": 3, "s_b_c": "qwe"}, "is_exixt": true}`


func main() {

	ret := new(interface{})
	jsonT.JsonLoads(jsonTestData, ret)

	retJson, _ := (*ret).(map[string]interface{})
	fmt.Println("Json get one ret is", retJson["one"])
	fmt.Println("Json get ext ret is", retJson["ext"])

	//jsonT.JsonIterGetOneFromByte()
	//
	//jsonT.JsonIterMashal()
	//
	//jsonT.JsonIterGet()

	jsonObj, err := jsonT.JsonIterLoad(jsonTestData)
	if nil != err {
		fmt.Println("JSON load fail err is", err)
		os.Exit(1)
	}

	ret1 := jsonObj("ext", "s_b").ToInt()
	fmt.Println("JSON get jsonTestData /ext/s_b/ ret is", ret1)
}