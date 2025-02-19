// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaerr_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/qioalice/ekago/v2/ekaerr"

	"github.com/stretchr/testify/assert"
)

type T struct{}

func (T) String() string {
	return "stringer"
}

func foo() *ekaerr.Error {
	return foo1().
		AddMessage("foo bad").
		AddFields("foo_arg", T{}).
		Throw()
}

func foo1() *ekaerr.Error {
	return foo2().
		AddMessage("foo1 bad").
		AddFields("foo1_arg", 23).
		Throw()
}

func foo2() *ekaerr.Error {
	return foo3().
		AddMessage("foo2 bad").
		AddFields("foo2_arg?", "").
		Throw()
}

func foo3() *ekaerr.Error {
	return ekaerr.IllegalState.
		New("what??", "arg1?", nil).
		Throw()
}

func TestError(t *testing.T) {
	foo().LogAsWarn()
}

func BenchmarkErrorCreation(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	defer func() {
		runtime.GC()
		fmt.Printf("EKAERR: %+v\n", ekaerr.EPS())
		fmt.Println()
	}()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = ekaerr.NotImplemented.New("An error")
	}
}

func BenchmarkErrorCreationReusing(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	defer func() {
		runtime.GC()
		fmt.Printf("EKAERR: %+v\n", ekaerr.EPS())
		fmt.Println()
	}()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := ekaerr.NotImplemented.New("An error")
		ekaerr.ReleaseError(&err)
	}
}

func BenchmarkErrorCreateAndReleaser(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := foo()
		ekaerr.ReleaseError(&err)
	}
}

func TestError_IsAnyDeep(t *testing.T) {
	cls := ekaerr.AlreadyExist.NewSubClass("Derived")
	err := cls.New("Error")

	assert.True(t, err.Is(cls))
	assert.False(t, err.Is(ekaerr.AlreadyExist))

	assert.True(t, err.IsAnyDeep(ekaerr.AlreadyExist))
	assert.False(t, err.IsAnyDeep(ekaerr.NotFound))
}
