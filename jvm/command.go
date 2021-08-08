package jvm

//lint:file-ignore ST1006 MYSTYLE
//lint:file-ignore U1000 MYSTYLE
// Gava虚拟机命令行参数解析
import (
	"flag"
	"fmt"
	"os"
)

// gava虚拟机命令行参数
const __HELP_FLAG_USAGE__ = "help will show the gava usage"
const __VERSION_FLAG_USAGE__ = "version will show the gava version"
const __CLASSPATH_FLAG_USAGE__ = "classpath will allow you to set gava virtual machine class path"

var SYS_JAVA_JRE_HOME string = ""

func init() {
	SYS_JAVA_JRE_HOME = os.Getenv("JAVA_HOME")
	if SYS_JAVA_JRE_HOME == "" {
		SYS_JAVA_JRE_HOME = os.Getenv("JRE_HOME")
	}
	debug("system java home => ", SYS_JAVA_JRE_HOME)
}

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
	flag.BoolVar(&command.Help, "help", false, __HELP_FLAG_USAGE__)
	flag.BoolVar(&command.Version, "version", false, __VERSION_FLAG_USAGE__)
	flag.StringVar(&command.ClassPath, "classpath", "", __CLASSPATH_FLAG_USAGE__)
	flag.StringVar(&command.ClassPath, "cp", "", __CLASSPATH_FLAG_USAGE__)
	flag.Parse()
	var args = flag.Args()
	if len(args) > 0 {
		command.EntryPointClass = args[0]
		command.Args = args[1:]
	}
	return command
}
