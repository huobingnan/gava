package jvm

import (
	"flag"
	"fmt"
)

// gava虚拟机命令行参数
const HELP_FLAG_USAGE = "help will show the gava usage"
const VERSION_FLAG_USAGE = "version will show the gava version"
const CLASSPATH_FLAG_USAGE = "classpath will allow you to set gava virtual machine class path"

type Command struct {
	Version         bool     // 是否显示版本号
	ClassPath       string   // classpath
	Help            bool     // 是否显示help
	EntryPointClass string   // 入口的class文件
	Args            []string // 运行时参数
}

// 解析gava虚拟机参数

func gavaUsage() {
	fmt.Println("usage: gava [options...] file [args..]")
}

func ParseCommand() Command {
	var command = Command{}
	flag.Usage = gavaUsage
	flag.BoolVar(&command.Help, "help", false, HELP_FLAG_USAGE)
	flag.BoolVar(&command.Version, "version", false, VERSION_FLAG_USAGE)
	flag.StringVar(&command.ClassPath, "classpath", "", CLASSPATH_FLAG_USAGE)
	flag.StringVar(&command.ClassPath, "cp", "", CLASSPATH_FLAG_USAGE)
	flag.Parse()
	var args = flag.Args()
	if len(args) > 0 {
		command.EntryPointClass = args[0]
		command.Args = args[1:]
	}
	return command
}
