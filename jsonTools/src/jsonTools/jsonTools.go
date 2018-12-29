package goTools

import (
	"encoding/json"
)


func JsonLoads(src string, dst interface{}) (error) {
	err := json.Unmarshal([]byte(src), dst)
	return err
}


func JsonDumps(src interface{}) (string, error) {
	var err 	error
	var dstByte []byte

	if dstByte, err = json.Marshal(src); err == nil {
		return string(dstByte), err
	}
	return "", err
}