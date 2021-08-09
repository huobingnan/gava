package runtime_test

import (
	"fmt"
	"gava/jvm"
	"testing"
)

func TestOperandStack(ctx *testing.T) {

	var stack = jvm.NewJvmOperandStack(10)

	stack.PushReference(&jvm.JObject{})
	stack.PushInt(10)
	stack.PushFloat(20.0)
	stack.PushDouble(30.0)
	stack.PushLong(40000)

	fmt.Println("PopLong => ", stack.PopLong() == 40000)
	fmt.Println("PopDouble => ", stack.PopDouble() == 30.0)
	fmt.Println("PopFloat => ", stack.PopFloat() == 20.0)
	fmt.Println("PopInt => ", stack.PopInt() == 10)
	fmt.Println("PopReference => ", stack.PopReference())

}
