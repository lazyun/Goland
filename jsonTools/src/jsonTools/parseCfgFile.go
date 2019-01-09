package jsonTools

import (
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"os"
)

func ParseCfgFile(fileName string, dst interface{}) {
	content, err := ioutil.ReadFile(fileName)
	if nil != err {
		panic( err )
	}

	if err = json.Unmarshal(content, dst); nil != err {
		panic( err )
	}
}


func GetCacheDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}
