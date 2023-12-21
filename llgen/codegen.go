// Copyright 2023 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package llgen

import (
	"langvm"
	"tinygo.org/x/go-llvm"
)

type Gen struct {
	LLModule llvm.Module
}

func (g *Gen) Dispose() {
	g.LLModule.Dispose()
}

func (g *Gen) Generate(mod langvm.Module) {
}

func (g *Gen) GenerateStruct() llvm.Type {
	var elementTypes []llvm.Type

	return llvm.StructType(elementTypes, false)
}
