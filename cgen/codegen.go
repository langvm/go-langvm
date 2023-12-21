// Copyright 2023 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package cgen

import (
	"fmt"
	"langvm"
	. "langvm/internal"
	"os"
)

type Gen struct {
	Types     []langvm.Type
	FuncImpls []*langvm.FuncImpl
}

func (g *Gen) GenerateType(typ langvm.Type) string {
	switch typ := typ.(type) {
	case langvm.BasicType:
		return g.GenerateBasicType(typ)
	case langvm.StructuredType:
		return g.GenerateStructType(typ)
	case langvm.FuncTrait:
		return g.GenerateTraitType(typ)
	case langvm.ArrayType:
		return g.GenerateArrayType(typ)
	}

	panic(fmt.Errorf("unexpected Type"))
}

func (g *Gen) GenerateBasicType(b langvm.BasicType) string {
	return map[langvm.BasicType]string{
		langvm.U8:  "uint8_t",
		langvm.U16: "uint16_t",
		langvm.U32: "uint32_t",
		langvm.U64: "uint64_t",
		langvm.S8:  "int8_t",
		langvm.S16: "int16_t",
		langvm.S32: "int32_t",
		langvm.S64: "int64_t",
		langvm.F32: "float",
		langvm.F64: "double",
	}[b]
}

func (g *Gen) GenerateGenDecl(typ langvm.Type, name string) string {
	switch typ := typ.(type) {
	case langvm.FuncType:
		return g.GenerateFuncPtrGenDecl(typ, name)
	case langvm.ArrayType:
		return g.GenerateArrayDecl(typ, name)
	case langvm.StructuredType:
		return g.GenerateType(typ) + " " + name
	case langvm.FuncTrait:
		return g.GenerateType(typ) + " " + name
	case langvm.BasicType:
		return g.GenerateType(typ) + " " + name
	}

	panic("unexpected Type")
}

func (g *Gen) GenerateStructType(structType langvm.StructuredType) (s string) {
	return "struct " + structType.Name
}

func (g *Gen) GenerateTraitType(traitType langvm.FuncTrait) string {
	return "struct " + traitType.Name
}

func (g *Gen) GenerateFuncResultType(funcType langvm.FuncType) (s string) {
	switch len(funcType.Results) {
	case 0:
		s += "void"
	case 1:
		s += g.GenerateType(funcType.Results[0].Type)
	default:
		s += "struct{"
		for _, r := range funcType.Results {
			s += g.GenerateGenDecl(r.Type, r.Name) + ";"
		}
		s += "}"
	}
	return
}

func (g *Gen) GenerateFuncDecl(funcType langvm.FuncType, name string) (s string) {

	s += g.GenerateFuncResultType(funcType) + " " + name + "("

	// T Name (T param1, T param2, ... T paramN)
	for i, p := range funcType.Params {
		s += g.GenerateGenDecl(p.Type, p.Name)
		if i != len(funcType.Params)-1 {
			s += ","
		}
	}

	s += ")"

	return
}

func (g *Gen) GenerateStructBody(structType langvm.StructuredType) (s string) {
	s = "struct " + structType.Name + "{"
	for _, field := range structType.Fields {
		s += g.GenerateGenDecl(field.Type, field.Name) + ";"
	}
	s += "}"

	return
}

func (g *Gen) GenerateTraitBody(traitType langvm.FuncTrait) (s string) {
	s += "struct " + traitType.Name + "{"
	for _, m := range traitType.Prototypes {
		s += g.GenerateFuncPtrGenDecl(m.Type.(langvm.FuncType), m.Name) + ";"
	}
	s += "}"
	return
}

func (g *Gen) GenerateFuncPtrGenDecl(funcType langvm.FuncType, name string) string {
	return g.GenerateFuncDecl(funcType, "(*"+name+")")
}

func (g *Gen) GenerateFuncType(funcType langvm.FuncType) string {
	return g.GenerateFuncPtrGenDecl(funcType, "")
}

func (g *Gen) GenerateFuncBody(impl *langvm.FuncImpl) (s string) {
	// TODO
	return
}

func (g *Gen) GenerateInvoke(impl *langvm.FuncImpl, invoke *langvm.Invoke) (s string) {
	if len(impl.FuncType.Params) != 0 {
		s += impl.Name + "_ret = "
	}
	s += impl.Name + "("
	// TODO
	s += ")"

	return
}

func (g *Gen) GenerateArrayDecl(arrayType langvm.ArrayType, name string) string {
	return g.GenerateType(arrayType.ElementType) + " " + name + "[" + Utoa(arrayType.Capacity, 16) + "]"
}

func (g *Gen) GenerateArrayType(arrayType langvm.ArrayType) string {
	return g.GenerateArrayDecl(arrayType, "")
}

func (g *Gen) GenerateHeader(out *os.File) (err error) {
	print := func(a ...any) {
		_, _ = fmt.Fprint(out, a...)
	}

	println := func(a ...any) {
		_, _ = fmt.Fprintln(out, a...)
	}

	println("#pragma once\n#include <stdint.h>")

	for _, typ := range g.Types {
		switch typ := typ.(type) {
		case langvm.StructuredType:
			print(g.GenerateStructType(typ), ";")
		case langvm.FuncTrait:
			print(g.GenerateTraitType(typ), ";")
		}
	}

	for _, typ := range g.Types {
		switch typ := typ.(type) {
		case langvm.StructuredType:
			print(g.GenerateStructBody(typ), ";")
		case langvm.FuncTrait:
			print(g.GenerateTraitBody(typ), ";")
		}
	}

	for _, f := range g.FuncImpls {
		print(g.GenerateFuncDecl(f.FuncType, f.Name), ";")
	}

	print("\n")

	return
}

func (g *Gen) GenerateSource(out *os.File) (err error) {
	err = g.GenerateHeader(out)
	if err != nil {
		return
	}

	return
}
