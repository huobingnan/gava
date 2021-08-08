package jvm

//lint:file-ignore ST1006 MYSTYLE
//lint:file-ignore U1000 MYSTYLE
import (
	"log"
	"os"
)

// 日志配置
const __DEBUG_ENABLE__ = true

var debug func(v ...interface{}) = func(v ...interface{}) {}

var info = log.New(os.Stdout, "[INFO] ", log.LstdFlags).Println
var fatal = log.New(os.Stderr, "[ERROR] ", log.LstdFlags).Fatal // 会终止程序运行

func init() {
	if __DEBUG_ENABLE__ {
		// 开启debug
		debug = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags).Println
		debug("DEBUG MODE ENABLED")
	}
}
