package goTools

import (
	"bytes"
	"encoding/gob"
	//"fmt"
	//"reflect"
)

/*
	此函文件主要处理了 gob 数据转换
	在 const 变量里定义需要注册的常量数值，在 RegistNumb 数组里添加需要注册的类型
 */

var RegistNumb = [...]interface{}{map[string]int{}, map[string]string{}}

const (
	RegistMSI = iota
	RegistMSS
)

var tempBuffer bytes.Buffer
var enc = gob.NewEncoder(&tempBuffer)
var dec = gob.NewDecoder(&tempBuffer)

/*
说明：将 map、array、silent、interface、struct 等类型转换为 byte 数据
输入：src 需要编码的数据，可选参数 index 需要注册的类型
输出：编码后的数据、可能的异常
*/
func FGobEncode(src interface{}, index ... int) (bytes.Buffer, error) {
	if index != nil {
		gob.Register(RegistNumb[index[0]])
	}

	err := enc.Encode(src)
	return tempBuffer, err
}


/*
说明：将 byte 数据转换为 map[string]int 的字典
输入：输出字典的指针
输出：可能的 error
*/
func FGobDecodeMSInt(dst *map[string]int) error {
	err := dec.Decode(dst)
	return err
}


/*
说明：将 byte 数据转换为 map[string]map[string]int 的字典
输入：输出字典的指针
输出：可能的 error
*/
func FGobDecodeMSSInt(dst *map[string]map[string]int) error {
	err := dec.Decode(dst)
	return err
}


/*
说明：将 byte 数据转换为 map[string]map[string]interface{} 的字典
输入：输出字典的指针
输出：可能的 error
*/
func FGobDecodeMSSItf(dst *map[string]map[string]interface{}) error {
	err := dec.Decode(dst)
	return err
}


func FGobDecodeMSSSItf(dst *map[string]map[string]map[string]interface{}) error {
	err := dec.Decode(dst)
	return err
}