package main

import (
	"fmt"
	"gava/jvm"
)

func main() {

	var command = jvm.ParseCommand()
	fmt.Println(command)
}
