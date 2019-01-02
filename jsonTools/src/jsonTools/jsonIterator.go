package jsonTools

import (
	"github.com/json-iterator/go"
	"fmt"
	//"reflect"
)

var TestData = ``

// **********************************************  test ***************************************************
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


func JsonIterGet(jsonString string, key_path ...interface{}) (jsoniter.Any) {
	val := []byte(jsonString)
	ret := jsoniter.Get(val, key_path ...)
	return ret
}


// **********************************************  usefull ***************************************************
func JsonIterLoad(jsonString string) (func(key_path ...interface{}) (jsoniter.Any) , error) {
	var jsonRet map[string]interface{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(jsonString), &jsonRet)
	if nil != err {
		return nil, err
	}

	varByte := []byte(jsonString)
	return func(key_path ...interface{}) (jsoniter.Any) {
		return jsoniter.Get(varByte, key_path ...)
	}, nil

}