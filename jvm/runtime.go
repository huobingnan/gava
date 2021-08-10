package jvm

import "math"

//lint:file-ignore ST1006 MYSTYLE
// JVM 运行时

// Java与Go类型的对照
type JByte = int8
type JShort = int16
type JChar = uint16
type JInt = int32
type JLong = int64
type JFloat = float32
type JDouble = float64
type JBoolean = bool

//#region Java Class Object
type JClass struct {
}

//#endregion

//#region Java Object
type JObject struct {
	class *JClass // 指向Class
}

//#endregion

//#region 栈帧
type JvmStackFrame struct {
	localVars    JvmLocalVars
	operandStack *JvmOperandStack
	next         *JvmStackFrame
	thread       *JvmThread
	nextPC       int
}

func NewJvmStackFrame(maxLocals uint, maxStack uint) *JvmStackFrame {
	return &JvmStackFrame{
		localVars:    NewJvmLocalVars(maxLocals),
		operandStack: NewJvmOperandStack(maxLocals),
		next:         nil,
	}
}
func (this *JvmStackFrame) OperandStack() *JvmOperandStack { return this.operandStack }
func (this *JvmStackFrame) LocalVars() JvmLocalVars        { return this.localVars }
func (this *JvmStackFrame) Next() *JvmStackFrame           { return this.next }
func (this *JvmStackFrame) Thread() *JvmThread             { return this.thread }
func (this *JvmStackFrame) NextPC() int                    { return this.nextPC }

type JvmSlot struct {
	number    int32    // 存放数字
	reference *JObject // 存放引用
}

type JvmLocalVars []JvmSlot

func NewJvmLocalVars(maxLocals uint) JvmLocalVars {
	if maxLocals > 0 {
		return make([]JvmSlot, maxLocals)
	} else {
		return nil
	}
}

func (this JvmLocalVars) SetInt(index uint, val int32) { this[index].number = val }

func (this JvmLocalVars) GetInt(index uint) int32 { return this[index].number }

func (this JvmLocalVars) SetFloat(index uint, val float32) {
	var bits = math.Float32bits(val)
	this[index].number = int32(bits)
}

func (this JvmLocalVars) GetFloat(index uint) float32 {
	return math.Float32frombits(uint32(this[index].number))
}

// long需要占据两个槽的空间
func (this JvmLocalVars) SetLong(index uint, val int64) {
	this[index].number = int32(val)         // 取后32位
	this[index+1].number = int32(val >> 32) // 去前32位
}

func (this JvmLocalVars) GetLong(index uint) int64 {
	var low = uint(this[index].number)
	var high = uint(this[index+1].number)
	return int64(high)<<32 | int64(low)
}

func (this JvmLocalVars) SetDouble(index uint, val float64) {
	var long = int64(math.Float64bits(val))
	this[index].number = int32(long)
	this[index+1].number = int32(long >> 32)
}

func (this JvmLocalVars) GetDouble(index uint) float64 {
	var low = uint32(this[index].number)
	var high = uint32(this[index+1].number)
	var bits = int64(high)<<32 | int64(low)
	return math.Float64frombits(uint64(bits))
}

func (this JvmLocalVars) SetReference(index uint, reference *JObject) {
	this[index].reference = reference
}

func (this JvmLocalVars) GetReference(index uint) *JObject { return this[index].reference }

type JvmOperandStack struct {
	top   uint
	slots []JvmSlot
}

func NewJvmOperandStack(maxStack uint) *JvmOperandStack {
	if maxStack > 0 {
		return &JvmOperandStack{
			top:   0,
			slots: make([]JvmSlot, maxStack),
		}
	} else {
		return nil
	}
}

// 在操作操作数栈时，不必考虑栈溢出和空栈问题，栈的大小在编译期已经确定了

func (this *JvmOperandStack) PushInt(val int32) {
	this.slots[this.top].number = val
	this.top++
}

func (this *JvmOperandStack) PopInt() int32 {
	this.top--
	return this.slots[this.top].number
}

func (this *JvmOperandStack) PushFloat(val float32) {
	this.slots[this.top].number = int32(math.Float32bits(val))
	this.top++
}

func (this *JvmOperandStack) PopFloat() float32 {
	this.top--
	return math.Float32frombits(uint32(this.slots[this.top].number))
}

func (this *JvmOperandStack) PushLong(val int64) {
	// long占两个槽的位置
	this.slots[this.top].number = int32(val) // 取低32位
	this.top++
	this.slots[this.top].number = int32(val >> 32) // 取高32位
	this.top++
}

func (this *JvmOperandStack) PopLong() int64 {
	// 取高32位（栈的特性）
	this.top--
	var high = uint32(this.slots[this.top].number)
	this.top--
	// 取低32位
	var low = uint32(this.slots[this.top].number)
	return int64(high)<<32 | int64(low)
}

func (this *JvmOperandStack) PushDouble(val float64) {
	var bits = math.Float64bits(val)
	this.PushLong(int64(bits))
}

func (this *JvmOperandStack) PopDouble() float64 {
	var long = this.PopLong()
	return math.Float64frombits(uint64(long))
}

func (this *JvmOperandStack) PushReference(val *JObject) {
	this.slots[this.top].reference = val
	this.top++
}

func (this *JvmOperandStack) PopReference() *JObject {
	this.top--
	var res = this.slots[this.top].reference
	this.slots[this.top].reference = nil // 使GC回收
	return res
}

func (this *JvmOperandStack) PushSlot(s JvmSlot) {
	this.slots[this.top] = s
	this.top++
}

func (this *JvmOperandStack) PopSlot() JvmSlot {
	this.top--
	return this.slots[this.top]
}

//#endregion

//#region 运行栈定义

type JvmStack struct {
	maxSize uint           // 虚拟机栈的最大容量
	size    uint           // 当前栈容量
	top     *JvmStackFrame // 栈顶
}

func NewJvmStack(maxSize uint) *JvmStack {

	return &JvmStack{
		maxSize: maxSize,
		size:    0,
		top:     nil,
	}
}

func (this *JvmStack) Push(frame *JvmStackFrame) {
	if this.size >= this.maxSize {
		panic("java.lang.StackOverflowError")
	}
	if this.top == nil {
		this.top = frame
	} else {
		frame.next = this.top
		this.top = frame
	}
	this.size++
}

func (this *JvmStack) Pop() *JvmStackFrame {
	if this.top == nil {
		panic("jvm stack empty error")
	}
	var res = this.top
	this.top = this.top.next
	return res
}

func (this *JvmStack) Peek() *JvmStackFrame { return this.top }

func (this *JvmStack) Size() uint { return this.size }

func (this *JvmStack) MaxSize() uint { return this.maxSize }

//endregion

//#region 运行时线程结构体定义

// JVM 线程定义
type JvmThread struct {
	pc    int       // 程序计数器
	stack *JvmStack // 运行时栈
}

func NewJvmThread() *JvmThread {
	return &JvmThread{
		pc:    -1,
		stack: NewJvmStack(1024),
	}
}

func (this *JvmThread) PC() int { return this.pc }

func (this *JvmThread) SetPC(pc int) { this.pc = pc }

func (this *JvmThread) PushFrame(frame *JvmStackFrame) { this.stack.Push(frame) }

func (this *JvmThread) PopFrame() *JvmStackFrame { return this.stack.Pop() }

func (this *JvmThread) CurrentFrame() *JvmStackFrame { return this.stack.Peek() }

//#endregion
