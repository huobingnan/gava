package main

import (
	"fmt"
	"gava/jvm"
)

func main() {

	var command = jvm.ParseCommand()
	fmt.Println(command)
	fmt.Println(jvm.SYS_JAVA_JRE_HOME)
}
