package goTools

import (
	"github.com/json-iterator/go"
	"fmt"
	//"reflect"
)

var TestData = ``

func JsonIterGetOneFromByte() {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	ret := jsoniter.Get(val, "Colors", 0).ToString()
	fmt.Println("JsonIterGetOneFromByte ret is", ret)
}


func JsonIterMashal() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	var jsonRet map[string]interface{}
	err := json.Unmarshal([]byte(TestData), &jsonRet)
	if nil != err {
		fmt.Println("Json Iterator marshal err",  err)
	}

	getValue := jsonRet["pkg_name"]
	fmt.Println("Json data get pkg_name value is", getValue )
}


func JsonIterGet() {
	val := []byte(TestData)
	ret := jsoniter.Get(val, "msg", "custom").ToString()
	fmt.Println("JsonIterGetOneFromByte get msg.custom ret is", ret)

	ret = jsoniter.Get(val, "msg").ToString()
	fmt.Println("JsonIterGetOneFromByte get msg ret is", ret)
}


func JsonIterLoad(jsonString string) (error, func(key_path ...interface{}) (jsoniter.Any)) {
	var jsonRet map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(TestData), &jsonRet)
	if nil != err {
		return err, nil
	}

	varByte := []byte(jsonString)
	return nil, func(key_path ...interface{}) (jsoniter.Any) {

		fmt.Println("emm", key_path)
		return jsoniter.Get(varByte, key_path ...)
	}

}