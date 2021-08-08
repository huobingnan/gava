package class_test

import (
	"fmt"
	"gava/jvm"
	"io/ioutil"
	"testing"
)

func TestReadClass(ctx *testing.T) {
	var bytecode, err = ioutil.ReadFile("/Users/huobingnan/tmp/gava-class-test/Hello.class")
	if err != nil {
		ctx.Fatal(err)
	}
	var javaClass *jvm.JavaClass
	javaClass, err = jvm.ParseJavaByteCode(bytecode)
	if err != nil {
		ctx.Fatal(err)
	}
	fmt.Printf("%v\n", javaClass)

}
