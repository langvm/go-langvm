// Copyright 2023 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package langvm

import (
	"unsafe"
)

const PointerSize = uint(unsafe.Sizeof(uintptr(0)))

const (
	TERMINATE = 0

	U_PTR = U_BEGIN + BasicType(PointerSize)
)

const (
	U_BEGIN BasicType = 010 + iota

	U8
	U16
	U32
	U64

	U_END
)
const (
	S_BEGIN BasicType = 020 + iota

	S8
	S16
	S32
	S64

	S_END
)

const (
	F_BEGIN BasicType = 030 + iota

	_
	_
	F32
	F64

	F_END
)

type TypeSignature struct {
	Sequence []BasicType
}

// Type is type.
type Type interface {
	Identify() string
	GetSignature() TypeSignature
}

type BasicType byte

func (b BasicType) Size() uint { return uint(b) & 07 }
func (b BasicType) Identify() string {
	return []string{
		U8:  "u8",
		U16: "u16",
		U32: "u32",
		U64: "u64",

		S8:  "i8",
		S16: "i16",
		S32: "i32",
		S64: "i64",

		F32: "f32",
		F64: "f64",
	}[b]
}
func (b BasicType) GetSignature() TypeSignature { return TypeSignature{[]BasicType{b}} }

func (b BasicType) IsUint() bool  { return b > U_BEGIN && b < U_END }
func (b BasicType) IsInt() bool   { return b > S_BEGIN && b < S_END }
func (b BasicType) IsFloat() bool { return b > F_BEGIN && b < F_END }

type (
	PointerType struct {
		PointeeType Type
	}

	UnsafePointerType struct {
		PointeeType Type
	}

	NullablePointerType struct {
		PointeeType Type

		// If true, the dangling pointer checker will **NOT** be applied.
		Checked bool
	}
)

func (p PointerType) Identify() string { return "ptrOf_" + p.PointeeType.Identify() }
func (p PointerType) GetSignature() TypeSignature {
	return TypeSignature{
		Sequence: append(append([]BasicType{U_PTR}, p.PointeeType.GetSignature().Sequence...), TERMINATE),
	}
}

func (p UnsafePointerType) Identify() string { return "ptrOf_" + p.PointeeType.Identify() }
func (p UnsafePointerType) GetSignature() TypeSignature {
	return TypeSignature{
		Sequence: append(append([]BasicType{U_PTR}, p.PointeeType.GetSignature().Sequence...), TERMINATE),
	}
}

func (p NullablePointerType) Identify() string { return "ptrOf_" + p.PointeeType.Identify() }
func (p NullablePointerType) GetSignature() TypeSignature {
	return TypeSignature{
		Sequence: append(append([]BasicType{U_PTR}, p.PointeeType.GetSignature().Sequence...), TERMINATE),
	}
}

type FuncType struct {
	Params, Results []Field
}

func (f FuncType) Identify() string { return "" } // TODO

type Field struct {
	Name string
	Type Type
}

// StructuredType is record-type.
// Recursive recording causes crashing so it is NOT allowed.
type StructuredType struct {
	Name          string
	Fields        []Field
	RecordedTypes map[string]bool

	Signature TypeSignature
}

func NewStructuredType() StructuredType {
	return StructuredType{
		Fields: []Field{},
	}
}

func (s StructuredType) Identify() string            { return s.Name }
func (s StructuredType) GetSignature() TypeSignature { return s.Signature }
func (s StructuredType) AddField(field Field) bool {
	if typ, ok := field.Type.(StructuredType); ok {
		for name := range typ.RecordedTypes {
			// Recursive record type.
			if name == s.Name {
				return false
			}
		}

		for name := range typ.RecordedTypes {
			s.RecordedTypes[name] = true
		}
	}

	s.Fields = append(s.Fields, field)
	s.Signature.Sequence = append(s.Signature.Sequence, field.Type.GetSignature().Sequence...)

	return true
}

type ArrayType struct {
	ElementType Type
	Capacity    uint
}

func (a ArrayType) Identify() string { return "" } // TODO
func (a ArrayType) GetSignature() TypeSignature {
	return TypeSignature{
		Sequence: append(append([]BasicType{U_PTR}, a.ElementType.GetSignature().Sequence...), TERMINATE),
	}
}

type TypeAlias struct {
	Name string
	Type Type
}

func (t TypeAlias) Identify() string            { return t.Name }
func (t TypeAlias) GetSignature() TypeSignature { return t.Type.GetSignature() }
