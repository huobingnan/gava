package runtime_test

import (
	"fmt"
	"gava/jvm"
	"testing"
)

func Test1(ctx *testing.T) {
	var localVars = jvm.NewJvmLocalVars(10)

	localVars.SetInt(0, 10)
	localVars.SetFloat(1, 30.3)
	localVars.SetDouble(2, 40.6)
	localVars.SetLong(4, 1000)
	localVars.SetReference(6, &jvm.JObject{})

	fmt.Println("GetInt => ", localVars.GetInt(0) == 10)
	fmt.Println("GetFloat => ", localVars.GetFloat(1) == 30.3)
	fmt.Println("GetDouble => ", localVars.GetDouble(2) == 40.6)
	fmt.Println("GetLong => ", localVars.GetLong(4) == 1000)
	fmt.Println("GetReference => ", localVars.GetReference(6))
}
