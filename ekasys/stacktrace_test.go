// Copyright © 2019. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GetStackTrace with 'skip' == 0, 'depth' == 1,
// tests:
// - GetStackTrace returns slice with len == 1 (as depth)
// - frame.Function contains current test name
func TestGetStackTraceCommonDepth1(t *testing.T) {

	frames := GetStackTrace(0, 1)

	assert.Len(t, frames, 1, "invalid len of frames")
	assert.Contains(t, frames[0].Function, "TestGetStackTraceCommonDepth1",
		"wrong function name")
}

// Test GetStackTrace with 'skip' == -3 (include hidden frames),
// 'depth' == -1 (full depth) tests:
// - GetStackTrace returns slice with len >= 3
// (at least hidden frames were included to the output)
// - first three returned frames have valid function names
func TestGetStackTraceCommonDepthAbsolutelyFull(t *testing.T) {

	frames := GetStackTrace(-3, -1)

	assert.True(t, len(frames) >= 3, "invalid len of frames")

	funcNames := []string{
		"runtime.Callers", "getStackFramePoints", "GetStackTrace",
	}

	for i := 0; i < len(funcNames) && i < len(frames); i++ {
		assert.Contains(t, frames[i].Function, funcNames[i],
			"wrong function name")
	}

	frames.Print(nil)
}

type T struct{}

func (T) foo() StackFrame {
	return GetStackTrace(0, 1)[0]
}

// TestStackFrame_DoFormat just see what StackFrame.DoFormat generates.
func TestStackFrame_DoFormat(t *testing.T) {

	frame := GetStackTrace(0, 1)[0]
	fmt.Println(frame.doFormat())

	frame = new(T).foo()
	fmt.Println(frame.doFormat())
}

// Bench StackFrame.doFormat func (generating readable output of stack frame).
func BenchmarkStackFrame_DoFormat(b *testing.B) {

	b.ReportAllocs()
	b.StopTimer()

	frame := GetStackTrace(0, 1)[0]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		frame.doFormat()
	}
}
