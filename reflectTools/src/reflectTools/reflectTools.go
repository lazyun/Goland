package reflectTools

import (
	"reflect"
	"fmt"
)


func ReflectTest(testInterface interface{}) {

	switch testInterface.(type) {
	case int:
		fmt.Println("this interface value is", testInterface, "type int ")
	case string:
		fmt.Println("this interface value is", testInterface, "type string")
	case float64:
		fmt.Println("this interface value is", testInterface, "type float64 ")
	case func():
		fmt.Println("this interface value is", testInterface, "type func ")
	}


	//value, ok := testInterface.(int)
	//fmt.Println("this interface value is", value, "type int ", ok)


}


func ReflectVarTypeValue(sm interface{}) {
	smType := reflect.TypeOf(sm)
	smValue := reflect.ValueOf(sm)

	fmt.Println("Receive type is", smType.String(), "value is", smValue.String())
}


func ReflectFuncDo(sm interface{}, values ...interface{}) {
	smValue := reflect.ValueOf(sm)
	if reflect.Func != smValue.Kind() {
		fmt.Println("This interface type is not Func!")
	}

	args := []reflect.Value{}

	for _, value := range values {
		args = append(args, reflect.ValueOf(value))
	}

	smValue.Call(args)
}