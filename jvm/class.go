package jvm

import (
	"encoding/binary"
	"fmt"
	"math"
	"unicode/utf16"
)

//lint:file-ignore ST1006  MYSTYLE
//lint:file-ignore U1000  MYSTYLE

// Java字节码读取
type JavaByteCodeReader struct {
	bytecode []byte
}

// 读取8位无符号整数 1 Byte
func (this *JavaByteCodeReader) ReadUint8() uint8 {
	var res = this.bytecode[0]
	// 利用切片的方式，左移bytecode
	this.bytecode = this.bytecode[1:]
	return res
}

// 读取16位无符号整数 2 Bytes
func (this *JavaByteCodeReader) ReadUint16() uint16 {
	var res = binary.BigEndian.Uint16(this.bytecode)
	this.bytecode = this.bytecode[2:]
	return res
}

// 读取32位无符号整数 4 Bytes
func (this *JavaByteCodeReader) ReadUint32() uint32 {
	var res = binary.BigEndian.Uint32(this.bytecode)
	this.bytecode = this.bytecode[4:]
	return res
}

// 读取64位无符号整数 8 Bytes
func (this *JavaByteCodeReader) ReadUint64() uint64 {
	var res = binary.BigEndian.Uint64(this.bytecode)
	this.bytecode = this.bytecode[8:]
	return res
}

// 读取uint16的数组，数组大小由开头数值决定
func (this *JavaByteCodeReader) ReadUint16s() []uint16 {
	var size = this.ReadUint16()
	var res = make([]uint16, size)
	for idx := range res {
		res[idx] = this.ReadUint16()
	}
	return res
}

// 读取制定大小字节的数据
func (this *JavaByteCodeReader) ReadBytes(size uint32) []byte {
	var res = this.bytecode[:size]
	this.bytecode = this.bytecode[size:]
	return res
}

// 常量池tag值定义
const (
	CONSTANT_Class              = 7
	CONSTANT_Fieldref           = 9
	CONSTANT_Methodref          = 10
	CONSTANT_InterfaceMethodref = 11
	CONSTANT_String             = 8
	CONSTANT_Integer            = 3
	CONSTANT_Float              = 4
	CONSTANT_Long               = 5
	CONSTANT_Double             = 6
	CONSTANT_NameAndType        = 12
	CONSTANT_Utf8               = 1
	CONSTANT_MethodHandle       = 15
	CONSTANT_MethodType         = 16
	CONSTANT_InvokeDynamic      = 18
)

// JVM预定义属性名称
const (
	CODE                 = "Code"
	CONSTANT_VALUE       = "ConstantValue"
	DEPRECATED           = "Deprecated"
	EXCEPTIONS           = "Exceptions"
	LINE_NUMBER_TABLE    = "LineNumberTable"
	LOCAL_VARIABLE_TABLE = "LocalVariableTable"
	SOURCE_FILE          = "SourceFile"
	SYNTHETIC            = "Synthetic"
)

// 常量信息接口定义
type ConstantInformation interface {
	ReadInformation(reader *JavaByteCodeReader)
}

// 派生常量信息

// 整型的常量信息（不仅仅用来存储int）
// 也用来存储比int更小的short，char，boolean，byte等
type ConstantIntegerInfo struct {
	intValue JInt
}

func (this *ConstantIntegerInfo) ReadInformation(reader *JavaByteCodeReader) {
	var val = reader.ReadUint32()
	this.intValue = JInt(val)
}

// 浮点数常量信息
type ConstantFloatInfo struct {
	floatValue JFloat
}

func (this *ConstantFloatInfo) ReadInformation(reader *JavaByteCodeReader) {
	var fbits = reader.ReadUint32()
	var float = math.Float32frombits(fbits)
	this.floatValue = JFloat(float)
}

// 长整形常量信息
type ConstantLongInfo struct {
	longValue JLong
}

func (this *ConstantLongInfo) ReadInformation(reader *JavaByteCodeReader) {
	var bits = reader.ReadUint64()
	this.longValue = JLong(bits)
}

// 双精度浮点型常量信息
type ConstantDoubleInfo struct {
	doubleValue JDouble
}

func (this *ConstantDoubleInfo) ReadInformation(reader *JavaByteCodeReader) {
	var bits = reader.ReadUint64()
	var double = math.Float64frombits(bits)
	this.doubleValue = JDouble(double)
}

// Utf8字符串常量信息. Class文件中的UTF8字串采用了MUTF8的编码格式，这里在解码时，需要使用MUTF8
type ConstantUtf8Info struct {
	stringValue string
}

// Copy 作者源码 => https://github.com/zxh0/jvmgo-book/blob/master/v1/code/go/src/jvmgo/ch03/classfile/cp_utf8.go
// 解析MUTF8
func __decodeMUtf8(bytearr []byte) string {
	utflen := len(bytearr)
	chararr := make([]uint16, utflen)

	var c, char2, char3 uint16
	count := 0
	chararr_count := 0

	for count < utflen {
		c = uint16(bytearr[count])
		if c > 127 {
			break
		}
		count++
		chararr[chararr_count] = c
		chararr_count++
	}

	for count < utflen {
		c = uint16(bytearr[count])
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			/* 0xxxxxxx*/
			count++
			chararr[chararr_count] = c
			chararr_count++
		case 12, 13:
			/* 110x xxxx   10xx xxxx*/
			count += 2
			if count > utflen {
				panic("malformed input: partial character at end")
			}
			char2 = uint16(bytearr[count-1])
			if char2&0xC0 != 0x80 {
				panic(fmt.Errorf("malformed input around byte %v", count))
			}
			chararr[chararr_count] = c&0x1F<<6 | char2&0x3F
			chararr_count++
		case 14:
			/* 1110 xxxx  10xx xxxx  10xx xxxx*/
			count += 3
			if count > utflen {
				panic("malformed input: partial character at end")
			}
			char2 = uint16(bytearr[count-2])
			char3 = uint16(bytearr[count-1])
			if char2&0xC0 != 0x80 || char3&0xC0 != 0x80 {
				panic(fmt.Errorf("malformed input around byte %v", (count - 1)))
			}
			chararr[chararr_count] = c&0x0F<<12 | char2&0x3F<<6 | char3&0x3F<<0
			chararr_count++
		default:
			/* 10xx xxxx,  1111 xxxx */
			panic(fmt.Errorf("malformed input around byte %v", count))
		}
	}
	// The number of chars produced may be less than utflen
	chararr = chararr[0:chararr_count]
	runes := utf16.Decode(chararr)
	return string(runes)
}

func (this *ConstantUtf8Info) ReadInformation(reader *JavaByteCodeReader) {
	var size = reader.ReadUint16() // 前16位标注了Utf8字串的长度信息
	var bytes = reader.ReadBytes(uint32(size))
	this.stringValue = __decodeMUtf8(bytes)
}

func (this *ConstantUtf8Info) String() string {
	return this.stringValue
}

// String常量信息， String常量本身不存放字符串信息，它指向了Utf8常量池
type ConstantStringInfo struct {
	cp          ConstantPool // 常量池
	stringIndex uint16       // 索引值
}

func (this *ConstantStringInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.stringIndex = reader.ReadUint16()
}

func (this *ConstantStringInfo) String() string {
	return this.cp.getUtf8(this.stringIndex)
}

// Class常量信息， 与StringConstant类似，Class的类信息字串也是存放于Utf8常量信息中
type ConstantClassInfo struct {
	cp        ConstantPool // 常量池
	nameIndex uint16       // 类信息索引
}

func (this *ConstantClassInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.nameIndex = reader.ReadUint16()
}

func (this *ConstantClassInfo) Name() string {
	return this.cp.getUtf8(this.nameIndex)
}

// 名称和类型描述
type ConstantNameAndTypeInfo struct {
	nameIndex       uint16 // 名称索引
	descriptorIndex uint16 // 描述索引
}

func (this *ConstantNameAndTypeInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.nameIndex = reader.ReadUint16()
	this.descriptorIndex = reader.ReadUint16()
}

// MemberrefInfo 常量信息
type ConstantMemberrefInfo struct {
	cp               ConstantPool // 常量池
	classIndex       uint16
	nameAndTypeIndex uint16
}

func (this *ConstantMemberrefInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.classIndex = reader.ReadUint16()
	this.nameAndTypeIndex = reader.ReadUint16()
}

func (this *ConstantMemberrefInfo) ClassName() string {
	return this.cp.getClassName(this.classIndex)
}

func (this *ConstantMemberrefInfo) NameAndDescriptor() (string, string) {
	return this.cp.getNameAndType(this.nameAndTypeIndex)
}

// Fieldref 常量信息
type ConstantFieldrefInfo struct{ ConstantMemberrefInfo }

func (this *ConstantFieldrefInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.classIndex = reader.ReadUint16()
	this.nameAndTypeIndex = reader.ReadUint16()
}

func (this *ConstantFieldrefInfo) ClassName() string { return this.cp.getClassName(this.classIndex) }

func (this *ConstantFieldrefInfo) NameAndDescriptor() (string, string) {
	return this.cp.getNameAndType(this.nameAndTypeIndex)
}

// Methodref 常量信息
type ConstantMethodrefInfo struct{ ConstantMemberrefInfo }

func (this *ConstantMethodrefInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.classIndex = reader.ReadUint16()
	this.nameAndTypeIndex = reader.ReadUint16()
}

func (this *ConstantMethodrefInfo) ClassName() string { return this.cp.getClassName(this.classIndex) }

func (this *ConstantMethodrefInfo) NameAndDescriptor() (string, string) {
	return this.cp.getNameAndType(this.nameAndTypeIndex)
}

// InterfaceMethodref 常量信息
type ConstantInterfaceMethodrefInfo struct{ ConstantMemberrefInfo }

func (this *ConstantInterfaceMethodrefInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.classIndex = reader.ReadUint16()
	this.nameAndTypeIndex = reader.ReadUint16()
}

func (this *ConstantInterfaceMethodrefInfo) ClassName() string {
	return this.cp.getClassName(this.classIndex)
}

func (this *ConstantInterfaceMethodrefInfo) NameAndDescriptor() (string, string) {
	return this.cp.getNameAndType(this.nameAndTypeIndex)
}

// MethodType 常量信息 JSE 1.7 引入
type ConstantMethodTypeInfo struct {
	descriptorIndex uint16
}

func (this *ConstantMethodTypeInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.descriptorIndex = reader.ReadUint16()
}

// ConstantMethodHandle 常量信息 JSE 1.7引入
type ConstantMethodHandleInfo struct {
	referenceKind  uint8
	referenceIndex uint16
}

func (this *ConstantMethodHandleInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.referenceKind = reader.ReadUint8()
	this.referenceIndex = reader.ReadUint16()
}

// ConstantInvokeDynamic 常量信息 JSE1.7 引入
type ConstantInvokeDynamicInfo struct {
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

func (this *ConstantInvokeDynamicInfo) ReadInformation(reader *JavaByteCodeReader) {
	this.bootstrapMethodAttrIndex = reader.ReadUint16()
	this.nameAndTypeIndex = reader.ReadUint16()
}

// 常量池
type ConstantPool struct {
	informations []ConstantInformation
}

//region 常量池私有函数

func (this *ConstantPool) getUtf8(stringIndex uint16) string {
	var utf8 = this.informations[stringIndex].(*ConstantUtf8Info)
	return utf8.String()
}

func (this *ConstantPool) getClassName(classIndex uint16) string {
	var class = this.informations[classIndex].(*ConstantClassInfo)
	return class.Name()
}

func (this *ConstantPool) getNameAndType(descriptorIndex uint16) (string, string) {
	var nameAndType = this.informations[descriptorIndex].(*ConstantNameAndTypeInfo)
	var name = this.getUtf8(nameAndType.nameIndex)
	var t = this.getUtf8(nameAndType.descriptorIndex)
	return name, t
}

//endregion

// 成员信息
type MemberInformation struct {
	cp              ConstantPool // 常量池
	accessFlags     uint16       // 访问标识符
	nameIndex       uint16
	descriptorIndex uint16
	attributes      []*Attribute
}

// 读取成员信息
func readMember(reader *JavaByteCodeReader, cp ConstantPool) *MemberInformation {
	// 成员信息依次为大端存储形式的
	//1. 2bytes => 访问标识符
	//2. 2bytes => 名称索引
	//3. 2bytes => 描述者索引
	return &MemberInformation{
		cp:              cp,
		accessFlags:     reader.ReadUint16(),
		nameIndex:       reader.ReadUint16(),
		descriptorIndex: reader.ReadUint16(),
		attributes:      readAttributes(reader, cp),
	}
}

//#region MemberInformation 读取常量池有关操作函数

// 属性信息
type Attribute interface {
	ReadAttribute(reader *JavaByteCodeReader)
}

// 这个属性十分重要, 顶层属性，可以套娃
type CodeAttribute struct {
	cp              ConstantPool
	name            string            // 属性名称
	length          uint32            // 属性值的长度
	maxStack        uint16            // 最大操作数栈深度
	maxLocals       uint16            // 局部变量表的最大值
	code            []byte            // 字节码
	exceptionTables []*ExceptionTable // 异常表
	attributes      []*Attribute      // 属性表
}

// 代码异常表
type ExceptionTable struct {
	startPC   uint16 // 程序计数器起始位置
	endPC     uint16 // 程序计数器结束位置
	handlerPC uint16 // 处理程序段起始位置
	catchType uint16 // 捕获的类型
}

func (this *CodeAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	this.maxStack = reader.ReadUint16()    // 首先是最大栈深度
	this.maxLocals = reader.ReadUint16()   // 最大变量表
	var codeSize = reader.ReadUint32()     // 字节码长度
	this.code = reader.ReadBytes(codeSize) // 读取字节码
	// 读取exceptions table
	var exceptionTableSize = reader.ReadUint16()
	this.exceptionTables = make([]*ExceptionTable, exceptionTableSize)
	for idx := range this.exceptionTables {
		this.exceptionTables[idx] = &ExceptionTable{
			startPC:   reader.ReadUint16(),
			endPC:     reader.ReadUint16(),
			handlerPC: reader.ReadUint16(),
			catchType: reader.ReadUint16(),
		}
	}
	this.attributes = readAttributes(reader, this.cp) // 套娃读取
}

type ConstantValueAttribute struct {
	name               string
	length             uint32
	constantValueIndex uint16
}

func (this *ConstantValueAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	this.constantValueIndex = reader.ReadUint16()
}

type DeprecatedAttribute struct {
	name   string
	length uint32
}

//do noting
func (this *DeprecatedAttribute) ReadAttribute(reader *JavaByteCodeReader) {}

// 异常属性
type ExceptionsAttribute struct {
	name                string
	length              uint32
	exceptionIndexTable []uint16
}

func (this *ExceptionsAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	this.exceptionIndexTable = reader.ReadUint16s()
}

func (this *ExceptionsAttribute) ExceptionsIndexTable() []uint16 { return this.exceptionIndexTable }

// 与异常处理有关
type LineNumberTableAttribute struct {
	name            string
	length          uint32
	lineNumberTable []*LineNumberTableEntry
}

type LineNumberTableEntry struct {
	startPC    uint16
	lineNumber uint16
}

func (this *LineNumberTableAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	var length = reader.ReadUint16()
	this.lineNumberTable = make([]*LineNumberTableEntry, length)
	for idx := range this.lineNumberTable {
		this.lineNumberTable[idx] = &LineNumberTableEntry{
			startPC:    reader.ReadUint16(),
			lineNumber: reader.ReadUint16(),
		}
	}
}

type LocalVariableTableAttribute struct {
	name               string
	length             uint32
	localVariableTable []*LocalVariableTableEntry
}

type LocalVariableTableEntry struct {
	startPc         uint16
	length          uint16
	nameIndex       uint16
	descriptorIndex uint16
	index           uint16
}

func (this *LocalVariableTableAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	var localVariableTableLength = reader.ReadUint16()
	this.localVariableTable = make([]*LocalVariableTableEntry, localVariableTableLength)
	for i := range this.localVariableTable {
		this.localVariableTable[i] = &LocalVariableTableEntry{
			startPc:         reader.ReadUint16(),
			length:          reader.ReadUint16(),
			nameIndex:       reader.ReadUint16(),
			descriptorIndex: reader.ReadUint16(),
			index:           reader.ReadUint16(),
		}
	}
}

type SourceFileAttribute struct {
	cp              ConstantPool
	name            string
	length          uint32
	sourceFileIndex uint16
}

func (this *SourceFileAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	this.sourceFileIndex = reader.ReadUint16()
}

func (this *SourceFileAttribute) FileName() string { return this.cp.getUtf8(this.sourceFileIndex) }

type SyntheticAttribute struct {
	name   string
	length uint32
}

// do noting
func (this *SyntheticAttribute) ReadAttribute(reader *JavaByteCodeReader) {}

type UnparsedAttribute struct {
	name        string
	length      uint32
	information []byte
}

func (this *UnparsedAttribute) ReadAttribute(reader *JavaByteCodeReader) {
	this.information = reader.ReadBytes(this.length)
}

// 读取属性信息表
func readAttributes(reader *JavaByteCodeReader, cp ConstantPool) []*Attribute {
	var attributeCount = reader.ReadUint16() // 2字节表示信息长度
	var attributes = make([]*Attribute, attributeCount)
	for idx := range attributes {
		// 构造属性信息表
		var attributeNameIndex = reader.ReadUint16()
		var attributeName = cp.getUtf8(attributeNameIndex)
		var attributeLength = reader.ReadUint32() // 信息长度有多少
		var attribute Attribute
		switch attributeName {
		case CODE:
			attribute = &CodeAttribute{cp: cp, name: attributeName, length: attributeLength}
		case CONSTANT_VALUE:
			attribute = &ConstantValueAttribute{name: attributeName, length: attributeLength}
		case DEPRECATED:
			attribute = &DeprecatedAttribute{name: attributeName, length: attributeLength}
		case EXCEPTIONS:
			attribute = &ExceptionsAttribute{name: attributeName, length: attributeLength}
		case LINE_NUMBER_TABLE:
			attribute = &LineNumberTableAttribute{name: attributeName, length: attributeLength}
		case LOCAL_VARIABLE_TABLE:
			attribute = &LocalVariableTableAttribute{name: attributeName, length: attributeLength}
		case SOURCE_FILE:
			attribute = &SourceFileAttribute{cp: cp, name: attributeName, length: attributeLength}
		case SYNTHETIC:
			attribute = &SyntheticAttribute{name: attributeName, length: attributeLength}
		default:
			attribute = &UnparsedAttribute{name: attributeName, length: attributeLength}
		}
		attribute.ReadAttribute(reader)
		attributes[idx] = &attribute
	}
	return attributes
}

// 读取所有的成员信息
func readMembers(reader *JavaByteCodeReader, cp ConstantPool) []*MemberInformation {
	var memberCount = reader.ReadUint16() // 首先读出成员的个数
	var members = make([]*MemberInformation, memberCount)
	for idx := range members {
		members[idx] = readMember(reader, cp)
	}
	return members
}

//#endregion

// Java Class对象
type JavaClass struct {
	magic          uint32               // Java字节码的魔数
	minorVersion   uint16               // 字节码副版本号
	majorVersion   uint16               // 字节码主版本号
	constantPool   ConstantPool         // 常量池
	accessFlags    uint16               // 访问标志符
	thisClass      uint16               // 当前Class
	superClass     uint16               // 超类
	interfaceClass []uint16             // 接口类
	fields         []*MemberInformation // 字段信息集合
	methods        []*MemberInformation // 方法信息集合
	attributes     []*Attribute         // 属性信息集合
}

//#region JavaClass getter & accessor

// 获取副版本号
func (this *JavaClass) MinorVersion() uint16 { return this.minorVersion }

// 获取主版本号
func (this *JavaClass) MajorVersion() uint16 { return this.majorVersion }

// 获取常量池
func (this *JavaClass) ConstantPool() ConstantPool { return this.constantPool }

// 获取访问标识符
func (this *JavaClass) AccessFlags() uint16 { return this.accessFlags }

// 获取字段信息
func (this *JavaClass) Fields() []*MemberInformation { return this.fields }

// 获取方法信息
func (this *JavaClass) Methods() []*MemberInformation { return this.methods }

// 获取类的全限定名称
func (this *JavaClass) ClassName() string { return "" }

// 获取超类的全限定名称
func (this *JavaClass) SuperClassName() string { return "" }

// 获取所有接口的全限定名称
func (this *JavaClass) InterfaceNames() []string { return []string{} }

//#endregion JavaClass getter & accessor

//#region JavaClass 解析字节码内容的私有方法，这些方法均可以安全的进行panic

// 读取字节码，构造Class对象
// 字节码的排列顺序
// 魔数 -> 次版本号 -> 主版本号 -> 常量池 -> 类访问标志 -> 两个uint16类型的常量池索引（本类和超类）
func (this *JavaClass) read(reader *JavaByteCodeReader) {
	this.readAndCheckMagicNumber(reader)                        // 检查魔数
	this.readAndCheckVersion(reader)                            // 检查版本号
	this.readConstantPool(reader)                               // 读取常量池信息
	this.accessFlags = reader.ReadUint16()                      // 读取类的访问标识符
	this.thisClass = reader.ReadUint16()                        // 本类
	this.superClass = reader.ReadUint16()                       // 超类
	this.interfaceClass = reader.ReadUint16s()                  // 读取接口信息
	this.fields = readMembers(reader, this.constantPool)        // 读取字段
	this.methods = readMembers(reader, this.constantPool)       // 读取方法
	this.attributes = readAttributes(reader, this.constantPool) // 读取属性

}

// 读取并检查魔数
func (this *JavaClass) readAndCheckMagicNumber(reader *JavaByteCodeReader) {
	var magic = reader.ReadUint32()
	if magic != 0xCAFEBABE {
		panic("java.lang.ClassFormatError => magic!")
	}
	this.magic = magic
}

// 读取并检查字节码文件版本号
func (this *JavaClass) readAndCheckVersion(reader *JavaByteCodeReader) {
	this.minorVersion = reader.ReadUint16() // 首先出现的是副版本号
	this.majorVersion = reader.ReadUint16() // 其次是主版本号

	debug("minor version => ", this.minorVersion)
	debug("major version => ", this.majorVersion)
	switch this.majorVersion {
	case 45:
		return
	case 46, 47, 48, 49, 50, 51, 52:
		if this.minorVersion == 0 {
			return
		}
	}
	panic("java.lang.UnsupportedClassVersionError")
}

// 读取解析常量池信息
func (this *JavaClass) readConstantPool(reader *JavaByteCodeReader) {
	var constantPoolSize = reader.ReadUint16()
	var constantPool = new(ConstantPool)
	var informations = make([]ConstantInformation, constantPoolSize)
	// 开始解析常量池
	// 一定注意，常量池的索引是从1开始的
	for i := 1; i < int(constantPoolSize); i++ {
		var tag = reader.ReadUint8() // 获取常量的tag
		var constantInformation ConstantInformation
		switch tag {
		case CONSTANT_Integer:
			constantInformation = &ConstantIntegerInfo{}
		case CONSTANT_Double:
			constantInformation = &ConstantDoubleInfo{}
		case CONSTANT_Float:
			constantInformation = &ConstantFloatInfo{}
		case CONSTANT_Long:
			constantInformation = &ConstantLongInfo{}
		case CONSTANT_Utf8:
			constantInformation = &ConstantUtf8Info{}
		case CONSTANT_String:
			constantInformation = &ConstantStringInfo{cp: *constantPool}
		case CONSTANT_Class:
			constantInformation = &ConstantClassInfo{cp: *constantPool}
		case CONSTANT_Fieldref:
			constantInformation = &ConstantFieldrefInfo{ConstantMemberrefInfo{cp: *constantPool}}
		case CONSTANT_Methodref:
			constantInformation = &ConstantMethodrefInfo{ConstantMemberrefInfo{cp: *constantPool}}
		case CONSTANT_InterfaceMethodref:
			constantInformation = &ConstantInterfaceMethodrefInfo{ConstantMemberrefInfo{cp: *constantPool}}
		case CONSTANT_NameAndType:
			constantInformation = &ConstantNameAndTypeInfo{}
		case CONSTANT_MethodType:
			constantInformation = &ConstantMethodTypeInfo{}
		case CONSTANT_MethodHandle:
			constantInformation = &ConstantMethodTypeInfo{}
		case CONSTANT_InvokeDynamic:
			constantInformation = &ConstantInvokeDynamicInfo{}
		default:
			panic("java.lang.ClassFormatError! => constant pool")

		}
		constantInformation.ReadInformation(reader) // 从字节码中读取信息
		informations[i] = constantInformation
		// http://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4.5
		// All 8-byte constants take up two entries in the constant_pool table of the class file.
		// If a CONSTANT_Long_info or CONSTANT_Double_info structure is the item in the constant_pool
		// table at index n, then the next usable item in the pool is located at index n+2.
		// The constant_pool index n+1 must be valid but is considered unusable.
		switch informations[i].(type) {
		case *ConstantLongInfo, *ConstantDoubleInfo:
			i++
		}
	}
	constantPool.informations = informations
	this.constantPool = *constantPool
}

//#endregion

// 解析Java字节码
func ParseJavaByteCode(bytecode []byte) (*JavaClass, error) {
	defer func() {
		var err = recover()
		if err != nil {
			fatal(err)
		}
	}()
	var reader = JavaByteCodeReader{bytecode: bytecode}
	var javaClass = JavaClass{}
	javaClass.read(&reader)
	return &javaClass, nil
}
