// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekasys

import (
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

// StackTrace is the slice of StackFrames, nothing more.
// Each stack level described separately.
type StackTrace []StackFrame

// getStackFramePoints returns the stack trace point's slice
// that contains 'count' points and starts from 'skip' depth level.
//
// You can pass any value <= 0 as 'count' to get full stack trace points.
func getStackFramePoints(skip, count int) (framePoints []uintptr) {

	// allow to get absolutely full stack trace
	// (include 'getStackFramePoints' and 'runtime.Callers' functions)
	if skip < -2 {
		skip = -2
	}

	// but by default (if 0 has been passed as 'skip', it means to skip
	// these functions ('getStackFramePoints' and 'runtime.Callers'),
	// Thus:
	// skip < -2 => skip = 0 => with these functions
	// skip == 0 => skip = 2 => w/o these functions
	skip += 2

	// get exactly as many stack trace frames as 'count' is
	// only if 'count' is more than zero
	if count > 0 {
		framePoints = make([]uintptr, count)
		return framePoints[:runtime.Callers(skip, framePoints)]
	}

	const (
		// how much frame points requested first time
		baseFullStackFramePointsLen int = 32

		// maximum requested frame points
		maxFullStackFramePointsLen int = 128
	)

	// runtime.Callers only fills slice we provide which
	// so, if slice is full, reallocate mem and try to request frames again
	framePointsLen := 0
	for count = baseFullStackFramePointsLen; ; count <<= 1 {

		framePoints = make([]uintptr, count)
		framePointsLen = runtime.Callers(skip, framePoints)

		if framePointsLen < count || count == maxFullStackFramePointsLen {
			break
		}
	}

	framePoints = framePoints[:framePointsLen]
	return framePoints[:len(framePoints)-1] // ignore Go internal functions
}

// GetStackTrace returns the stack trace as StackFrame object's slice,
// that have specified 'depth' and starts from 'skip' depth level.
// Each StackFrame object represents an one stack trace depth-level.
//
// You can pass any value <= 0 as 'depth' to get full stack trace.
func GetStackTrace(skip, depth int) (stacktrace StackTrace) {

	// see the same code section in 'getStackFramePoints'
	// to more details what happening here with 'skip' arg
	if skip < -3 {
		skip = -3
	}
	skip++

	// prepare to get runtime.Frame objects:
	// - get stack trace frame points,
	// - create runtime.Frame iterator by frame points from prev step
	framePoints := getStackFramePoints(skip, depth)
	framePointsLen := len(framePoints)
	frameIterator := runtime.CallersFrames(framePoints)

	// alloc mem for slice that will have as many 'runtime.Frame' objects
	// as many frame points we got
	stacktrace = make([]StackFrame, framePointsLen)

	// i in func scope (not in loop's) because it will be 'frames' len
	i := 0
	for more := true; more && i < framePointsLen; i++ {
		stacktrace[i].Frame, more = frameIterator.Next()
	}

	// but frameIterator can provide less 'runtime.Frame' objects
	// than we requested -> should fix 'frames' len w/o reallocate
	return stacktrace[:i]
}

var (
	cachedStackFrames = make(map[uintptr]StackFrame)
	cachedStackFramesMu sync.RWMutex
)

func GetStackTrace2(skip, depth int) (stacktrace StackTrace) {

	// see the same code section in 'getStackFramePoints'
	// to more details what happening here with 'skip' arg
	if skip < -3 {
		skip = -3
	}
	skip++

	framePoints := getStackFramePoints(skip, depth)
	stacktrace = make(StackTrace, len(framePoints))

	needToUpdateCache := false
	found := false

	cachedStackFramesMu.RLock()
	for i, n := 0, len(framePoints); i < n; i++ {
		if stacktrace[i], found = cachedStackFrames[framePoints[i]]; !found {
			needToUpdateCache = true
		}
	}
	cachedStackFramesMu.RUnlock()

	if !needToUpdateCache {
		return
	}

	for i, n := 0, len(framePoints); i < n; i++ {
		if stacktrace[i].PC == 0 {
			fp := framePoints[i]

			f := runtime.FuncForPC(fp)
			file, line := f.FileLine(fp)

			stacktrace[i] = StackFrame{
				Frame: runtime.Frame{

					// Keep PC not initialized (0) because it's used as a signal
					// that StackFrame has been initialized just now
					// in the loop of updating cached frames below.
					//
					//PC:       framePoints[i],

					Func:     f,
					Function: f.Name(),
					File:     file,
					Line:     line,
					Entry:    f.Entry(),
				},
			}

			stacktrace[i].doFormat()
		}
	}

	cachedStackFramesMu.Lock()
	for i, n := 0, len(framePoints); i < n; i++ {
		fp := framePoints[i]
		if stacktrace[i].PC == 0 {
			stacktrace[i].PC = fp
		}
		cachedStackFrames[fp] = stacktrace[i]
	}
	cachedStackFramesMu.Unlock()

	return
}

// ExcludeInternal returns stacktrace based on current but with excluded all
// Golang internal functions such as runtime.doInit, runtime.main, etc.
func (s StackTrace) ExcludeInternal() StackTrace {

	// because some internal golang functions (such as runtime.gopanic)
	// could be embedded to user's function stacktrace's part,
	// we can't just cut and drop last part of stacktrace when we found
	// a function with a "runtime." prefix from the beginning to end.
	// instead, we starting from the end and generating "ignore list" -
	// a special list of stack frames that won't be included to the result set.

	idx := len(s) - 1
	for continue_ := true; continue_; idx-- {
		continue_ = idx > 0 && (strings.HasPrefix(s[idx].Function, "runtime.") ||
			strings.HasPrefix(s[idx].Function, "testing."))
	}
	return s[:idx+2]
}

// Write writes generated stacktrace to the w or to the stdout if w == nil.
func (s StackTrace) Write(w io.Writer) (n int, err error) {

	if w == nil {
		w = os.Stdout
	}

	for _, frame := range s {

		nn, err_ := w.Write([]byte(frame.DoFormat()))
		if err_ != nil {
			return n, err_
		}
		n += nn

		// write \n
		if _, err_ := w.Write([]byte{'\n'}); err_ != nil {
			return n, err_
		}
		n += 1
	}

	return n, nil
}

// Print prints generated stacktrace to the w or to the stdout if w == nil.
// Ignores all errors. To write with error tracking use Write method.
func (s StackTrace) Print(w io.Writer) {
	_, _ = s.Write(w)
}
