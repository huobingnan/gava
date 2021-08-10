package jvm

import "math"

//lint:file-ignore ST1006 MY

type InstructionCodeReader struct {
	bytecode []byte
	pc       int
}

func NewInstructionCodeReader(code []byte, pc int) *InstructionCodeReader {
	return &InstructionCodeReader{
		bytecode: code,
		pc:       pc,
	}
}

func (this *InstructionCodeReader) Reset(code []byte, pc int) {
	this.bytecode = code
	this.pc = pc
}

func (this *InstructionCodeReader) ReadUint8() uint8 {
	var res = this.bytecode[this.pc]
	this.pc++
	return res
}

func (this *InstructionCodeReader) ReadInt8() int8 {
	return int8(this.ReadUint8())
}

func (this *InstructionCodeReader) ReadUint16() uint16 {
	var byte1 = uint16(this.ReadUint8())
	var byte2 = uint16(this.ReadUint8())
	return byte1<<8 | byte2
}

func (this *InstructionCodeReader) ReadInt16() int16 {
	return int16(this.ReadUint16())
}

func (this *InstructionCodeReader) ReadInt32() int32 {
	var byte1 = int32(this.ReadUint8())
	var byte2 = int32(this.ReadUint8())
	var byte3 = int32(this.ReadUint8())
	var byte4 = int32(this.ReadUint8())
	return byte1<<24 | byte2<<16 | byte3<<8 | byte4
}

type Instruction interface {
	FetchOperands(reader *InstructionCodeReader)
	Execute(frame *JvmStackFrame)
}

type NoOperandsInstruction struct{}

func (this *NoOperandsInstruction) FetchOperands(reader *InstructionCodeReader) {}

func (this *NoOperandsInstruction) Execute(fram *JvmStackFrame) {}

type BranchInstruction struct {
	Offset int
}

func (this *BranchInstruction) FetchOperands(reader *InstructionCodeReader) {
	this.Offset = int(reader.ReadUint16())
}

type Index8Instruction struct {
	Index uint
}

func (this *Index8Instruction) FetchOperands(reader *InstructionCodeReader) {
	this.Index = uint(reader.ReadUint8())
}

type Index16Instruction struct {
	Index uint
}

func (this *Index16Instruction) FetchOperands(reader *InstructionCodeReader) {
	this.Index = uint(reader.ReadUint16())
}

// CONST => 直接在操作数栈中压入一个数据
type NOP struct{ NoOperandsInstruction }
type ACONST_NULL struct{ NoOperandsInstruction }
type DCONST_0 struct{ NoOperandsInstruction }
type DCONST_1 struct{ NoOperandsInstruction }
type FCONST_0 struct{ NoOperandsInstruction }
type FCONST_1 struct{ NoOperandsInstruction }
type FCONST_2 struct{ NoOperandsInstruction }
type ICONST_M1 struct{ NoOperandsInstruction }
type ICONST_0 struct{ NoOperandsInstruction }
type ICONST_1 struct{ NoOperandsInstruction }
type ICONST_2 struct{ NoOperandsInstruction }
type ICONST_3 struct{ NoOperandsInstruction }
type ICONST_4 struct{ NoOperandsInstruction }
type ICONST_5 struct{ NoOperandsInstruction }
type LCONST_0 struct{ NoOperandsInstruction }
type LCONST_1 struct{ NoOperandsInstruction }

func (this *NOP) Execute(frame *JvmStackFrame)         {}
func (this *ACONST_NULL) Execute(frame *JvmStackFrame) { frame.operandStack.PushReference(nil) }
func (this *DCONST_0) Execute(frame *JvmStackFrame)    { frame.operandStack.PushDouble(0.0) }
func (this *DCONST_1) Execute(frame *JvmStackFrame)    { frame.operandStack.PushDouble(1.0) }
func (this *FCONST_0) Execute(frame *JvmStackFrame)    { frame.operandStack.PushFloat(0.0) }
func (this *FCONST_1) Execute(frame *JvmStackFrame)    { frame.operandStack.PushFloat(1.0) }
func (this *FCONST_2) Execute(frame *JvmStackFrame)    { frame.operandStack.PushFloat(2.0) }
func (this *ICONST_M1) Execute(frame *JvmStackFrame)   { frame.operandStack.PushInt(-1) }
func (this *ICONST_0) Execute(frame *JvmStackFrame)    { frame.operandStack.PushInt(0) }
func (this *ICONST_1) Execute(frame *JvmStackFrame)    { frame.operandStack.PushInt(1) }
func (this *ICONST_2) Execute(frame *JvmStackFrame)    { frame.operandStack.PushInt(2) }
func (this *ICONST_3) Execute(frame *JvmStackFrame)    { frame.operandStack.PushInt(3) }
func (this *ICONST_4) Execute(frame *JvmStackFrame)    { frame.operandStack.PushInt(4) }
func (this *ICONST_5) Execute(frame *JvmStackFrame)    { frame.operandStack.PushInt(5) }
func (this *LCONST_0) Execute(frame *JvmStackFrame)    { frame.operandStack.PushLong(0) }
func (this *LCONST_1) Execute(frame *JvmStackFrame)    { frame.operandStack.PushLong(1) }

// 读取一个字节转换为整数类型，然后推入栈顶
type BIPUSH struct{ value int8 }

func (this *BIPUSH) FetchOperands(reader *InstructionCodeReader) { this.value = reader.ReadInt8() }
func (this *BIPUSH) Execute(frame *JvmStackFrame)                { frame.operandStack.PushInt(int32(this.value)) }

// 读取两个字节，扩展为short，然后推入栈顶
type SIPUSH struct{ value int16 }

func (this *SIPUSH) FetchOperands(reader *InstructionCodeReader) { this.value = reader.ReadInt16() }
func (this *SIPUSH) Execute(frame *JvmStackFrame)                { frame.operandStack.PushInt(int32(this.value)) }

// LOAD => 从局部变量表取数，推送至操作数栈栈顶
type ILOAD struct{ Index8Instruction }
type LLOAD struct{ Index8Instruction }
type FLOAD struct{ Index8Instruction }
type DLOAD struct{ Index8Instruction }
type ALOAD struct{ Index8Instruction }
type ILOAD_0 struct{ NoOperandsInstruction }
type ILOAD_1 struct{ NoOperandsInstruction }
type ILOAD_2 struct{ NoOperandsInstruction }
type ILOAD_3 struct{ NoOperandsInstruction }
type LLOAD_0 struct{ NoOperandsInstruction }
type LLOAD_1 struct{ NoOperandsInstruction }
type LLOAD_2 struct{ NoOperandsInstruction }
type LLOAD_3 struct{ NoOperandsInstruction }
type FLOAD_0 struct{ NoOperandsInstruction }
type FLOAD_1 struct{ NoOperandsInstruction }
type FLOAD_2 struct{ NoOperandsInstruction }
type FLOAD_3 struct{ NoOperandsInstruction }
type DLOAD_0 struct{ NoOperandsInstruction }
type DLOAD_1 struct{ NoOperandsInstruction }
type DLOAD_2 struct{ NoOperandsInstruction }
type DLOAD_3 struct{ NoOperandsInstruction }
type ALOAD_0 struct{ NoOperandsInstruction }
type ALOAD_1 struct{ NoOperandsInstruction }
type ALOAD_2 struct{ NoOperandsInstruction }
type ALOAD_3 struct{ NoOperandsInstruction }

func __genericIloadImpl(frame *JvmStackFrame, index uint) {
	var value = frame.localVars.GetInt(index)
	frame.operandStack.PushInt(value)
}

func __genericLloadImpl(frame *JvmStackFrame, index uint) {
	var value = frame.localVars.GetLong(index)
	frame.operandStack.PushLong(value)
}

func __genericFloadImpl(frame *JvmStackFrame, index uint) {
	var value = frame.localVars.GetFloat(index)
	frame.operandStack.PushFloat(value)
}

func __genericDloadImpl(frame *JvmStackFrame, index uint) {
	var value = frame.localVars.GetDouble(index)
	frame.operandStack.PushDouble(value)
}

func __genericAloadImpl(frame *JvmStackFrame, index uint) {
	var value = frame.localVars.GetReference(index)
	frame.operandStack.PushReference(value)
}

func (this *ILOAD) Execute(frame *JvmStackFrame)   { __genericIloadImpl(frame, this.Index) }
func (this *LLOAD) Execute(frame *JvmStackFrame)   { __genericLloadImpl(frame, this.Index) }
func (this *FLOAD) Execute(frame *JvmStackFrame)   { __genericLloadImpl(frame, this.Index) }
func (this *DLOAD) Execute(frame *JvmStackFrame)   { __genericLloadImpl(frame, this.Index) }
func (this *ALOAD) Execute(frame *JvmStackFrame)   { __genericLloadImpl(frame, this.Index) }
func (this *ILOAD_0) Execute(frame *JvmStackFrame) { __genericIloadImpl(frame, 0) }
func (this *ILOAD_1) Execute(frame *JvmStackFrame) { __genericIloadImpl(frame, 1) }
func (this *ILOAD_2) Execute(frame *JvmStackFrame) { __genericIloadImpl(frame, 2) }
func (this *ILOAD_3) Execute(frame *JvmStackFrame) { __genericIloadImpl(frame, 3) }
func (this *LLOAD_0) Execute(frame *JvmStackFrame) { __genericLloadImpl(frame, 0) }
func (this *LLOAD_1) Execute(frame *JvmStackFrame) { __genericLloadImpl(frame, 1) }
func (this *LLOAD_2) Execute(frame *JvmStackFrame) { __genericLloadImpl(frame, 2) }
func (this *LLOAD_3) Execute(frame *JvmStackFrame) { __genericLloadImpl(frame, 3) }
func (this *FLOAD_0) Execute(frame *JvmStackFrame) { __genericFloadImpl(frame, 0) }
func (this *FLOAD_1) Execute(frame *JvmStackFrame) { __genericFloadImpl(frame, 1) }
func (this *FLOAD_2) Execute(frame *JvmStackFrame) { __genericFloadImpl(frame, 2) }
func (this *FLOAD_3) Execute(frame *JvmStackFrame) { __genericFloadImpl(frame, 3) }
func (this *DLOAD_0) Execute(frame *JvmStackFrame) { __genericDloadImpl(frame, 0) }
func (this *DLOAD_1) Execute(frame *JvmStackFrame) { __genericDloadImpl(frame, 1) }
func (this *DLOAD_2) Execute(frame *JvmStackFrame) { __genericDloadImpl(frame, 2) }
func (this *DLOAD_3) Execute(frame *JvmStackFrame) { __genericDloadImpl(frame, 3) }
func (this *ALOAD_0) Execute(frame *JvmStackFrame) { __genericAloadImpl(frame, 0) }
func (this *ALOAD_1) Execute(frame *JvmStackFrame) { __genericAloadImpl(frame, 1) }
func (this *ALOAD_2) Execute(frame *JvmStackFrame) { __genericAloadImpl(frame, 2) }
func (this *ALOAD_3) Execute(frame *JvmStackFrame) { __genericAloadImpl(frame, 3) }

// STORE => 从栈顶取数放到本地变量表
type ISTORE struct{ Index8Instruction }
type LSTORE struct{ Index8Instruction }
type FSTORE struct{ Index8Instruction }
type DSTORE struct{ Index8Instruction }
type ASTORE struct{ Index8Instruction }
type ISTORE_0 struct{ NoOperandsInstruction }
type ISTORE_1 struct{ NoOperandsInstruction }
type ISTORE_2 struct{ NoOperandsInstruction }
type ISTORE_3 struct{ NoOperandsInstruction }
type LSTORE_0 struct{ NoOperandsInstruction }
type LSTORE_1 struct{ NoOperandsInstruction }
type LSTORE_2 struct{ NoOperandsInstruction }
type LSTORE_3 struct{ NoOperandsInstruction }
type FSTORE_0 struct{ NoOperandsInstruction }
type FSTORE_1 struct{ NoOperandsInstruction }
type FSTORE_2 struct{ NoOperandsInstruction }
type FSTORE_3 struct{ NoOperandsInstruction }
type DSTORE_0 struct{ NoOperandsInstruction }
type DSTORE_1 struct{ NoOperandsInstruction }
type DSTORE_2 struct{ NoOperandsInstruction }
type DSTORE_3 struct{ NoOperandsInstruction }
type ASTORE_0 struct{ NoOperandsInstruction }
type ASTORE_1 struct{ NoOperandsInstruction }
type ASTORE_2 struct{ NoOperandsInstruction }
type ASTORE_3 struct{ NoOperandsInstruction }

func __genericIStoreImpl(frame *JvmStackFrame, index uint) {
	var value = frame.operandStack.PopInt()
	frame.localVars.SetInt(index, value)
}

func __genericLStoreImpl(frame *JvmStackFrame, index uint) {
	var value = frame.operandStack.PopLong()
	frame.localVars.SetLong(index, value)
}

func __genericFStoreImpl(frame *JvmStackFrame, index uint) {
	var value = frame.operandStack.PopFloat()
	frame.localVars.SetFloat(index, value)
}

func __genericDStoreImpl(frame *JvmStackFrame, index uint) {
	var value = frame.operandStack.PopDouble()
	frame.localVars.SetDouble(index, value)
}

func __genericAStoreImpl(frame *JvmStackFrame, index uint) {
	var value = frame.operandStack.PopReference()
	frame.localVars.SetReference(index, value)
}

func (this *ISTORE) Execute(frame *JvmStackFrame)   { __genericIStoreImpl(frame, this.Index) }
func (this *LSTORE) Execute(frame *JvmStackFrame)   { __genericLStoreImpl(frame, this.Index) }
func (this *FSTORE) Execute(frame *JvmStackFrame)   { __genericLStoreImpl(frame, this.Index) }
func (this *DSTORE) Execute(frame *JvmStackFrame)   { __genericLStoreImpl(frame, this.Index) }
func (this *ASTORE) Execute(frame *JvmStackFrame)   { __genericLStoreImpl(frame, this.Index) }
func (this *ISTORE_0) Execute(frame *JvmStackFrame) { __genericIStoreImpl(frame, 0) }
func (this *ISTORE_1) Execute(frame *JvmStackFrame) { __genericIStoreImpl(frame, 1) }
func (this *ISTORE_2) Execute(frame *JvmStackFrame) { __genericIStoreImpl(frame, 2) }
func (this *ISTORE_3) Execute(frame *JvmStackFrame) { __genericIStoreImpl(frame, 3) }
func (this *LSTORE_0) Execute(frame *JvmStackFrame) { __genericLStoreImpl(frame, 0) }
func (this *LSTORE_1) Execute(frame *JvmStackFrame) { __genericLStoreImpl(frame, 1) }
func (this *LSTORE_2) Execute(frame *JvmStackFrame) { __genericLStoreImpl(frame, 2) }
func (this *LSTORE_3) Execute(frame *JvmStackFrame) { __genericLStoreImpl(frame, 3) }
func (this *FSTORE_0) Execute(frame *JvmStackFrame) { __genericFStoreImpl(frame, 0) }
func (this *FSTORE_1) Execute(frame *JvmStackFrame) { __genericFStoreImpl(frame, 1) }
func (this *FSTORE_2) Execute(frame *JvmStackFrame) { __genericFStoreImpl(frame, 2) }
func (this *FSTORE_3) Execute(frame *JvmStackFrame) { __genericFStoreImpl(frame, 3) }
func (this *DSTORE_0) Execute(frame *JvmStackFrame) { __genericDStoreImpl(frame, 0) }
func (this *DSTORE_1) Execute(frame *JvmStackFrame) { __genericDStoreImpl(frame, 1) }
func (this *DSTORE_2) Execute(frame *JvmStackFrame) { __genericDStoreImpl(frame, 2) }
func (this *DSTORE_3) Execute(frame *JvmStackFrame) { __genericDStoreImpl(frame, 3) }
func (this *ASTORE_0) Execute(frame *JvmStackFrame) { __genericAStoreImpl(frame, 0) }
func (this *ASTORE_1) Execute(frame *JvmStackFrame) { __genericAStoreImpl(frame, 1) }
func (this *ASTORE_2) Execute(frame *JvmStackFrame) { __genericAStoreImpl(frame, 2) }
func (this *ASTORE_3) Execute(frame *JvmStackFrame) { __genericAStoreImpl(frame, 3) }

// 和操作数栈相关的指令
type POP struct{ NoOperandsInstruction }     // 弹出栈顶
type POP2 struct{ NoOperandsInstruction }    // 连续两次弹出
type DUP struct{ NoOperandsInstruction }     // 复制栈顶元素并压入栈顶
type DUP_X1 struct{ NoOperandsInstruction }  // 复制栈顶数值并将两个复制值压入栈顶
type DUP_X2 struct{ NoOperandsInstruction }  // 复制栈顶数值并将三个（或两个）复制值压入栈顶
type DUP2 struct{ NoOperandsInstruction }    // 复制栈顶一个（对于 long 或 double 类型）或两个数值（对于非 long 或 double 的其他类型）并将复制值压入栈顶
type DUP2_X1 struct{ NoOperandsInstruction } // dup_x1 指令的双倍版本
type DUP2_X2 struct{ NoOperandsInstruction } // dup_x2 指令的双倍版本
type SWAP struct{ NoOperandsInstruction }    // 交换栈顶两个变量

func (this *POP) Execute(frame *JvmStackFrame) { frame.operandStack.PopSlot() }

func (this *POP2) Execute(frame *JvmStackFrame) {
	frame.operandStack.PopSlot()
	frame.operandStack.PopSlot()
}

func (this *DUP) Execute(frame *JvmStackFrame) {
	var top = frame.operandStack.PopSlot()
	frame.operandStack.PushSlot(top)
	frame.operandStack.PushSlot(top)
}

/*
bottom -> top
[...][c][b][a]
          __/
         |
         V
[...][c][a][b][a]
*/
func (this *DUP_X1) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var slot1 = stack.PopSlot()
	var slot2 = stack.PopSlot()
	stack.PushSlot(slot1)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

/*
bottom -> top
[...][c][b][a]
       _____/
      |
      V
[...][a][c][b][a]
*/
func (this *DUP_X2) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var slot1 = stack.PopSlot()
	var slot2 = stack.PopSlot()
	var slot3 = stack.PopSlot()
	stack.PushSlot(slot1)
	stack.PushSlot(slot3)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

/*
bottom -> top
[...][c][b][a]____
          \____   |
               |  |
               V  V
[...][c][b][a][b][a]
*/
func (this *DUP2) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var slot1 = stack.PopSlot()
	var slot2 = stack.PopSlot()
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

/*
bottom -> top
[...][c][b][a]
       _/ __/
      |  |
      V  V
[...][b][a][c][b][a]
*/
func (this *DUP2_X1) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var slot1 = stack.PopSlot()
	var slot2 = stack.PopSlot()
	var slot3 = stack.PopSlot()
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
	stack.PushSlot(slot3)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

/*
bottom -> top
[...][d][c][b][a]
       ____/ __/
      |   __/
      V  V
[...][b][a][d][c][b][a]
*/
func (this *DUP2_X2) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var slot1 = stack.PopSlot()
	var slot2 = stack.PopSlot()
	var slot3 = stack.PopSlot()
	var slot4 = stack.PopSlot()
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
	stack.PushSlot(slot4)
	stack.PushSlot(slot3)
	stack.PushSlot(slot2)
	stack.PushSlot(slot1)
}

func (this *SWAP) Execute(frame *JvmStackFrame) {
	var slot1 = frame.operandStack.PopSlot()
	var slot2 = frame.operandStack.PopSlot()
	frame.operandStack.PushSlot(slot1)
	frame.operandStack.PushSlot(slot2)
}

// JVM数学指令实现
// 加法指令
type IADD struct{ NoOperandsInstruction }
type LADD struct{ NoOperandsInstruction }
type FADD struct{ NoOperandsInstruction }
type DADD struct{ NoOperandsInstruction }

// 减法指令
type ISUB struct{ NoOperandsInstruction }
type LSUB struct{ NoOperandsInstruction }
type FSUB struct{ NoOperandsInstruction }
type DSUB struct{ NoOperandsInstruction }

// 乘法指令
type IMUL struct{ NoOperandsInstruction }
type LMUL struct{ NoOperandsInstruction }
type FMUL struct{ NoOperandsInstruction }
type DMUL struct{ NoOperandsInstruction }

// 除法指令
type IDIV struct{ NoOperandsInstruction }
type LDIV struct{ NoOperandsInstruction }
type FDIV struct{ NoOperandsInstruction }
type DDIV struct{ NoOperandsInstruction }

// 求余数指令
type IREM struct{ NoOperandsInstruction }
type LREM struct{ NoOperandsInstruction }
type FREM struct{ NoOperandsInstruction }
type DREM struct{ NoOperandsInstruction }

//取反指令
type INEG struct{ NoOperandsInstruction }
type LNEG struct{ NoOperandsInstruction }
type FNEG struct{ NoOperandsInstruction }
type DNEG struct{ NoOperandsInstruction }

// 加法指令实现

func (this *IADD) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopInt()
	var i2 = frame.operandStack.PopInt()
	var add = i1 + i2
	frame.operandStack.PushInt(add)
}

func (this *LADD) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopLong()
	var i2 = frame.operandStack.PopLong()
	var add = i1 + i2
	frame.operandStack.PushLong(add)
}

func (this *FADD) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopFloat()
	var i2 = frame.operandStack.PopFloat()
	var add = i1 + i2
	frame.operandStack.PushFloat(add)
}

func (this *DADD) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopDouble()
	var i2 = frame.operandStack.PopDouble()
	var add = i1 + i2
	frame.operandStack.PushDouble(add)
}

// 减法指令实现

func (this *ISUB) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopInt()
	var i2 = frame.operandStack.PopInt()
	var sub = i2 - i1
	frame.operandStack.PushInt(sub)
}

func (this *LSUB) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopLong()
	var i2 = frame.operandStack.PopLong()
	var sub = i2 - i1
	frame.operandStack.PushLong(sub)
}

func (this *DSUB) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopDouble()
	var i2 = frame.operandStack.PopDouble()
	var sub = i2 - i1
	frame.operandStack.PushDouble(sub)
}

func (this *FSUB) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopFloat()
	var i2 = frame.operandStack.PopFloat()
	var sub = i2 - i1
	frame.operandStack.PushFloat(sub)
}

// 乘法指令实现

func (this *IMUL) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopInt()
	var i2 = frame.operandStack.PopInt()
	var res = i2 * i1
	frame.operandStack.PushInt(res)
}

func (this *LMUL) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopLong()
	var i2 = frame.operandStack.PopLong()
	var res = i2 * i1
	frame.operandStack.PushLong(res)
}

func (this *DMUL) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopDouble()
	var i2 = frame.operandStack.PopDouble()
	var res = i2 * i1
	frame.operandStack.PushDouble(res)
}

func (this *FMUL) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopFloat()
	var i2 = frame.operandStack.PopFloat()
	var res = i2 * i1
	frame.operandStack.PushFloat(res)
}

// 除法指令

func (this *IDIV) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopInt()
	var i2 = frame.operandStack.PopInt()
	var res = i2 / i1
	if i1 == 0 {
		panic("java.lang.ArithmeticException: division by zero")
	}
	frame.operandStack.PushInt(res)
}

func (this *LDIV) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopLong()
	var i2 = frame.operandStack.PopLong()
	var res = i2 / i1
	if i1 == 0 {
		panic("java.lang.ArithmeticException: division by zero")
	}
	frame.operandStack.PushLong(res)
}

func (this *DDIV) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopDouble()
	var i2 = frame.operandStack.PopDouble()
	var res = i2 / i1
	frame.operandStack.PushDouble(res)
}

func (this *FDIV) Execute(frame *JvmStackFrame) {
	var i1 = frame.operandStack.PopFloat()
	var i2 = frame.operandStack.PopFloat()
	var res = i2 / i1
	frame.operandStack.PushFloat(res)
}

// 求余运算

func (this *DREM) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var v2 = stack.PopDouble()
	var v1 = stack.PopDouble()
	var result = math.Mod(v1, v2)
	stack.PushDouble(result)
}

func (this *FREM) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var v2 = float64(stack.PopFloat())
	var v1 = float64(stack.PopFloat())
	var result = float32(math.Mod(v1, v2))
	stack.PushFloat(result)
}

func (this *IREM) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var v2 = stack.PopInt()
	var v1 = stack.PopInt()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: division by zero")
	}
	var result = v1 % v2
	stack.PushInt(result)
}

func (this *LREM) Execute(frame *JvmStackFrame) {
	var stack = frame.OperandStack()
	var v2 = stack.PopLong()
	var v1 = stack.PopLong()
	if v2 == 0 {
		panic("java.lang.ArithmeticException: division by zero")
	}
	var result = v1 % v2
	stack.PushLong(result)
}
