// Copyright © 2021. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"testing"
)

// benchGetStackTraceCommonDepth aux bench func that starts
// GetStackTrace bench with specified 'depth' arg.
func benchGetStackTraceCommonDepth(b *testing.B, depth int) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetStackTrace(0, depth)
	}
}

// benchGetStackTrace2CommonDepth aux bench func that starts
// GetStackTrace2 bench with specified 'depth' arg.
func benchGetStackTrace2CommonDepth(b *testing.B, depth int) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetStackTrace2(0, depth)
	}
}

// benchGetStackTraceSyntheticDepth increases stack depth level artificially
// by 'createDepth' value and then starts GetStackTrace bench with
// specified 'depth' arg.
func benchGetStackTraceSyntheticDepth(b *testing.B, depth, createDepth int) {

	type tF func(int, int) StackTrace

	wrapper := func(f tF) tF {
		return func(i1 int, i2 int) StackTrace {
			return f(i1, i2)
		}
	}

	var f tF = GetStackTrace

	for i := 0; i < createDepth; i++ {
		f = wrapper(f)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f(0, depth)
	}
}

// benchGetStackTrace2SyntheticDepth increases stack depth level artificially
// by 'createDepth' value and then starts GetStackTrace2 bench with
// specified 'depth' arg.
func benchGetStackTrace2SyntheticDepth(b *testing.B, depth, createDepth int) {

	type tF func(int, int) StackTrace

	wrapper := func(f tF) tF {
		return func(i1 int, i2 int) StackTrace {
			return f(i1, i2)
		}
	}

	var f tF = GetStackTrace2

	for i := 0; i < createDepth; i++ {
		f = wrapper(f)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f(0, depth)
	}
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 1 on common stack.
func Benchmark_GetStackTrace_CommonDepth_1(b *testing.B) {
	benchGetStackTraceCommonDepth(b, 1)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 10 on common stack.
func Benchmark_GetStackTrace_CommonDepth_10(b *testing.B) {
	benchGetStackTraceCommonDepth(b, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on common stack.
func Benchmark_GetStackTrace_CommonDepth_Full(b *testing.B) {
	benchGetStackTraceCommonDepth(b, -1)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 1
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_1_of_10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 1, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_10_of_10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 10, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_Full_of_10(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, -1, 10)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_10_of_20(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, 10, 20)
}

// Bench GetStackTrace with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace_SyntheticDepth_Full_of_20(b *testing.B) {
	benchGetStackTraceSyntheticDepth(b, -1, 20)
}

// ---------------------------------------------------------------------------- //

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == 1 on common stack.
func Benchmark_GetStackTrace2_CommonDepth_1(b *testing.B) {
	benchGetStackTrace2CommonDepth(b, 1)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == 10 on common stack.
func Benchmark_GetStackTrace2_CommonDepth_10(b *testing.B) {
	benchGetStackTrace2CommonDepth(b, 10)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == -1 (full depth)
// on common stack.
func Benchmark_GetStackTrace2_CommonDepth_Full(b *testing.B) {
	benchGetStackTrace2CommonDepth(b, -1)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == 1
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace2_SyntheticDepth_1_of_10(b *testing.B) {
	benchGetStackTrace2SyntheticDepth(b, 1, 10)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace2_SyntheticDepth_10_of_10(b *testing.B) {
	benchGetStackTrace2SyntheticDepth(b, 10, 10)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace2_SyntheticDepth_Full_of_10(b *testing.B) {
	benchGetStackTrace2SyntheticDepth(b, -1, 10)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == 10
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace2_SyntheticDepth_10_of_20(b *testing.B) {
	benchGetStackTrace2SyntheticDepth(b, 10, 20)
}

// Bench GetStackTrace2 with 'skip' == 0, 'depth' == -1 (full depth)
// on artificially enlarged stack by 10.
func Benchmark_GetStackTrace2_SyntheticDepth_Full_of_20(b *testing.B) {
	benchGetStackTrace2SyntheticDepth(b, -1, 20)
}
