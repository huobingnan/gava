package jvm

//lint:file-ignore ST1006 MYSTYLE
// JVM classpath 处理模块，针对输入的classpath进行解析。
// 并通过实现ClassEntry接口，实现对classpath下的文件进行读取

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const __OS_PATH_SEPARATOR__ = string(os.PathListSeparator) // 系统路径分隔符

// 日志配置
//var info = log.New(os.Stdout, "[jvm/classpath.go] ", log.LstdFlags).Println
var fatal = log.New(os.Stderr, "[jvm/classpath.go] ", log.LstdFlags).Fatal

type ClassEntry interface {
	ReadClass(classQulifierName string) ([]byte, ClassEntry, error)
	String() string
}

//+-------------------------------- CompositeClassEntry definition ---------------------------+

type CompositeClassEntry struct {
	entrys             []ClassEntry
	entrysAbsolutePath []string
	classpath          string // 对应的classpath
}

func (this *CompositeClassEntry) ReadClass(classQulifierName string) ([]byte, ClassEntry, error) {
	return nil, nil, nil
}

func (this *CompositeClassEntry) String() string {
	return this.classpath
}

func newCompositeClassEntry(classpath string) *CompositeClassEntry {
	var entrysPath = strings.Split(classpath, __OS_PATH_SEPARATOR__)
	if len(entrysPath) <= 1 {
		fatal("invalid system path separator")
	}
	// 构建绝对路径
	var entrysAbsolutePath = make([]string, len(entrysPath))
	for idx, path := range entrysPath {
		var abs, err = filepath.Abs(path)
		if err != nil {
			fatal("invalid classpath => ", path)
		}
		// 检测文件是否存在
		_, err = os.Stat(abs)
		if err != nil {
			fatal("classpath is not exists => ", path)
		}
		entrysAbsolutePath[idx] = abs
	}
	// 解析组合起来的classEntry
	var entrys = make([]ClassEntry, len(entrysPath))
	for idx, path := range entrysPath {
		entrys[idx] = NewClassEntry(path)
	}
	return &CompositeClassEntry{entrysAbsolutePath: entrysAbsolutePath, entrys: entrys, classpath: classpath}
}

//+---------------------------------- DirClassEntry definition ------------------------------------+

type DirClassEntry struct {
	entrysAbsolutePath string
	classpath          string
}

func (this *DirClassEntry) ReadClass(classQulifierName string) ([]byte, ClassEntry, error) {
	return nil, nil, nil
}

func (this *DirClassEntry) String() string {
	return this.entrysAbsolutePath
}

func newDirClassEntry(classpath string) *DirClassEntry {
	var absPath, err = filepath.Abs(classpath)
	if err != nil {
		fatal("invalid classpath => ", classpath)
	}
	// 检测文件是否存在
	var stat fs.FileInfo
	stat, err = os.Stat(absPath)
	if err != nil {
		fatal("classpath is not exists => ", classpath)
	}
	// 检测是否是dir
	if !stat.IsDir() {
		fatal("classpath is not a valid dir => ", classpath)
	}
	return &DirClassEntry{entrysAbsolutePath: absPath, classpath: classpath}
}

//+--------------------------------- CompressedClassEntry -------------------------------------+

type CompressedClassEntry struct {
	entrysAbsolutePath string
	compressedType     string
	classpath          string
}

func (this *CompressedClassEntry) ReadClass(classQulifierName string) ([]byte, ClassEntry, error) {
	return nil, nil, nil
}

func (this *CompressedClassEntry) String() string {
	return this.entrysAbsolutePath
}

func newCompressedClassEntry(classpath string) *CompressedClassEntry {
	var absPath, err = filepath.Abs(classpath)
	if err != nil {
		fatal("invalid classpath => ", classpath)
	}
	// 检测文件是否存在
	_, err = os.Stat(absPath)
	if err != nil {
		fatal("classpath is not exists => ", classpath)
	}

	var spliteResult = strings.Split(classpath, ".")
	var compressedType = spliteResult[len(spliteResult)-1]

	return &CompressedClassEntry{entrysAbsolutePath: absPath, classpath: classpath, compressedType: compressedType}
}

//+--------------------------------- Package functions ------------------------------------+

func NewClassEntry(classpath string) ClassEntry {
	if strings.Contains(classpath, __OS_PATH_SEPARATOR__) {
		return newCompositeClassEntry(classpath)
	} else if strings.HasSuffix(classpath, "jar") || strings.HasSuffix(classpath, "zip") {
		return newCompressedClassEntry(classpath)
	} else {
		return newDirClassEntry(classpath)
	}
}
