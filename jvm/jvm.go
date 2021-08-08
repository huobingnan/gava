package jvm

//lint:file-ignore ST1006 MYSTYLE
//lint:file-ignore U1000 MYSTYLE
import (
	"flag"
	"fmt"
	"log"
	"os"
)

// 日志配置
const __DEBUG_ENABLE__ = true

var debug func(v ...interface{})

var info = log.New(os.Stdout, "[INFO] ", log.LstdFlags).Println
var fatal = log.New(os.Stderr, "[ERROR] ", log.LstdFlags).Fatal // 会终止程序运行

func init() {
	if __DEBUG_ENABLE__ {
		debug = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags).Println
	} else {
		debug = func(v ...interface{}) {}
	}
}

// gava虚拟机命令行参数
const __HELP_FLAG_USAGE__ = "help will show the gava usage"
const __VERSION_FLAG_USAGE__ = "version will show the gava version"
const __CLASSPATH_FLAG_USAGE__ = "classpath will allow you to set gava virtual machine class path"

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
