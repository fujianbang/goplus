package exec

import (
	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------
// Push/Pop

type iPush struct {
	v interface{}
}

type iPop int
type iPopN int

func (p *iPush) Exec(stk *Stack, ctx *Context) {
	stk.Push(p.v)
}

func (p iPop) Exec(stk *Stack, ctx *Context) {
	stk.Pop()
}

func (p iPopN) Exec(stk *Stack, ctx *Context) {
	stk.data = stk.data[:len(stk.data)-int(p)]
}

func Push(v interface{}) Instr {
	return &iPush{v}
}

func PopN(n int) Instr {
	return iPopN(n)
}

var (
	Nil Instr = Push(nil)
	Pop Instr = iPop(0)
)

// -----------------------------------------------------------------------------
// Compose

type iCompose int

func (p iCompose) Exec(stk *Stack, ctx *Context) {

	n := len(stk.data)
	n1 := n-int(p)
	stk.data[n1] = stk.data[n-1]
	stk.data = stk.data[:n1+1]
}

func Compose(n int) Instr {
	return iCompose(n)
}

// -----------------------------------------------------------------------------
// Clear

type iClear int

func (p iClear) Exec(stk *Stack, ctx *Context) {
	stk.data = stk.data[:ctx.base]
}

var (
	Clear Instr = iClear(0)
)

// -----------------------------------------------------------------------------
// Or/And

type iOr int
type iAnd int

func (delta iOr) Exec(stk *Stack, ctx *Context) {
	a, _ := stk.Pop()
	if a1, ok := a.(bool); ok {
		if a1 {
			stk.Push(true)
			ctx.ip += int(delta)
		}
	} else {
		panic("left operand of || operator isn't a boolean expression")
	}
}

func (delta iAnd) Exec(stk *Stack, ctx *Context) {
	a, _ := stk.Pop()
	if a1, ok := a.(bool); ok {
		if !a1 {
			stk.Push(false)
			ctx.ip += int(delta)
		}
	} else {
		panic("left operand of && operator isn't a boolean expression")
	}
}

func Or(delta int) Instr {
	return iOr(delta)
}

func And(delta int) Instr {
	return iAnd(delta) 
}

// -----------------------------------------------------------------------------
// Jmp

type iJmp int

func (delta iJmp) Exec(stk *Stack, ctx *Context) {
	ctx.ip +=int(delta)
}

func Jmp(delta int) Instr {
	return iJmp(delta)
}

// -----------------------------------------------------------------------------
// JmpIfFalse

type iJmpIfFalse int

func (delta iJmpIfFalse) Exec(stk *Stack, ctx *Context) {
	a, _ := stk.Pop()
	if a1, ok := a.(bool); ok {
		if !a1 {
			ctx.ip += int(delta)
		}
	} else {
		panic("condition isn't a boolean expression")
	}
}

func JmpIfFalse(delta int) Instr {
	return iJmpIfFalse(delta)
}

// -----------------------------------------------------------------------------
// Case/Default

type iCase int

func (delta iCase) Exec(stk *Stack, ctx *Context) {
	b, _ := stk.Pop()
	a, _ := stk.Top()
	cond := qlang.EQ(a, b)
	if cond1, ok := cond.(bool); ok {
		if cond1 {
			stk.Pop()
		} else {
			ctx.ip += int(delta)
		}
	} else {
		panic("operator == return non-boolean value?")
	}
}

func Case(delta int) Instr {
	return iCase(delta)
}

var (
	Default Instr = Pop
)

// -----------------------------------------------------------------------------
// SubSlice

type iOp3 struct {
	op    func(v, a, b interface{}) interface{}
	arity int
	hasA  bool
	hasB  bool
}

func (p *iOp3) Exec(stk *Stack, ctx *Context) {
	var i = 1
	var a, b interface{}
	args := stk.PopNArgs(p.arity)
	if p.hasA {
		a = args[i]
		i++
	}
	if p.hasB {
		b = args[i]
	}
	stk.Push(p.op(args[0], a, b))
}

func Op3(op func(v, a, b interface{}) interface{}, hasA, hasB bool) Instr {
	n := 1
	if hasA { n++ }
	if hasB { n++ }
	return &iOp3{op, n, hasA, hasB}
}

// -----------------------------------------------------------------------------
