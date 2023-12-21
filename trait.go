// Copyright 2023 LangVM Project
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package langvm

type FuncPrototype struct {
	Name string
	Type FuncType
}

type FuncTrait struct {
	Name       string
	Prototypes []FuncPrototype
}

func NewFuncTrait() FuncTrait {
	return FuncTrait{
		Prototypes: []FuncPrototype{},
	}
}

func (t FuncTrait) Identify() string { return t.Name }
func (t FuncTrait) GetSignature() TypeSignature {
	return TypeSignature{
		Sequence: []BasicType{}, // TODO
	}
}

type OperationTrait struct {
}
