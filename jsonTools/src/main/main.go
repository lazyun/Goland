package main

import (
	jsonT "../jsonTools"
	"fmt"
)


func main() {
	jsonT.JsonIterGetOneFromByte()

	jsonT.JsonIterMashal()

	jsonT.JsonIterGet()

	err, jsonObj := jsonT.JsonIterLoad(jsonT.TestData)
	if nil != err {
		fmt.Println("JSON load fail err is", err)
	}

	ret := jsonObj("ext_info", "battery_infos", "level").ToInt()
	fmt.Println("JSON get /ext_info/battery_infos/level ret is", ret)
}