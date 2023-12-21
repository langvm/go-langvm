// Copyright 2023 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package langvm

type (
	Callee interface {
		Identify() string
		Type() FuncType
	}

	// NativeFunc defines the external function implemented in C.
	NativeFunc struct {
		Name string
		FuncType

		Symbol string
	}

	// FuncInline is inline and can only be inline.
	// Recursion inline invoke is not allowed in inline functions.
	FuncInline struct {
		Name string
		FuncType

		BasicBlocks []BasicBlock
	}

	// FuncImpl is a real function.
	FuncImpl struct {
		Name string
		FuncType

		BasicBlocks []BasicBlock
	}
)

func (f NativeFunc) Identify() string { return f.Name }
func (f NativeFunc) Type() FuncType   { return f.FuncType }

func (f FuncInline) Identify() string { return f.Name }
func (f FuncInline) Type() FuncType   { return f.FuncType }

func (f FuncImpl) Identify() string { return f.Name }
func (f FuncImpl) Type() FuncType   { return f.FuncType }

const (
	COND_EQ = iota
	COND_NEQ
	COND_LTEQ
	COND_LT
)

const (
	OP_ADD = iota
	OP_SUB
	OP_MUL
	OP_DIV
)

type (
	BasicBlock struct{}

	Invoke struct{}

	Branch struct {
		To uint // Basic block index.

		Cond     uint
		Operands [2]uint
	}

	ArithmeticOp struct {
		Op     uint
		Save   uint
		Source [2]uint
	}

	Jump struct {
		To uint // Basic block index.
	}
)
