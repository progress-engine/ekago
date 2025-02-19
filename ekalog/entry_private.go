// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekalog

import (
	"runtime"

	"github.com/qioalice/ekago/v2/ekasys"
	"github.com/qioalice/ekago/v2/internal/ekafield"
	"github.com/qioalice/ekago/v2/internal/ekaletter"
)

// prepare prepares current Entry for being used assuming that Entry has been
// obtained from the Error's pool. Returns prepared Entry.
func (e *Entry) prepare() *Entry {

	e.LogLetter.Items.Flags = DefaultFlags

	// Because we can't detect when work with Entry is done while user chains
	// Logger's method (with Entry cloning), we must use runtime.SetFinalizer()
	// to return Entry to the its pool.
	if e.needSetFinalizer {
		runtime.SetFinalizer(e, releaseEntryForFinalizer)
		e.needSetFinalizer = false
	}

	return e
}

// cleanup frees all allocated resources (RAM in 99% cases) by Entry e, preparing
// it for being returned to the pool and being reused in the future.
func (e *Entry) cleanup() (this *Entry) {

	e.l = nil
	e.LogLetter.StackTrace = nil
	e.ErrLetter = nil
	ekaletter.ResetItem(e.LogLetter.Items)

	return e
}

// clone clones the current Entry 'e' and returns it copy. It takes a new Entry
// object from its pool to avoid unnecessary RAM allocations.
func (e *Entry) clone() *Entry {

	clonedEntry := acquireEntry()

	clonedEntry.LogLetter.Items.Flags = e.LogLetter.Items.Flags

	// Clone Fields using most efficient way.
	// Do not allocate RAM if it's already allocated (but nulled).
	if lFrom := len(e.LogLetter.Items.Fields); lFrom > 0 {
		if cTo := cap(clonedEntry.LogLetter.Items.Fields); cTo < lFrom {
			clonedEntry.LogLetter.Items.Fields = make([]ekafield.Field, lFrom)
		} else {
			// lFrom <= cTo, it's ok to do that
			clonedEntry.LogLetter.Items.Fields =
				clonedEntry.LogLetter.Items.Fields[:lFrom]
		}
		for i := 0; i < lFrom; i++ {
			clonedEntry.LogLetter.Items.Fields[i] = e.LogLetter.Items.Fields[i]
		}
	}

	// There is no need to zero Time, Level, Message fields
	// because they used only in one place and will be overwritten anyway.

	return clonedEntry
}

// addFields extract key-value pairs from 'args' and adds it to the e's *LetterItem
// or just saving 'explicitFields. Returns this.
//
// Requirements:
// 'e' != nil. Otherwise UB (may panic).
func (e *Entry) addFields(args []interface{}, explicitFields []ekafield.Field) *Entry {

	// REMINDER!
	// IT IS STRONGLY GUARANTEES THAT BOTH OF 'args' AND 'explicitFields'
	// CAN NOT BE AN EMPTY (OR SET) AT THE SAME TIME!

	ekaletter.ParseTo(e.LogLetter.Items, args, explicitFields, true)
	return e
}

//func (e *Entry) forceStacktrace(ignoreFrames int) (this *Entry) {
//
//	e.setFlag(bEntryFlagAutoGenerateCaller)
//	e.ssf = ignoreFrames
//	return e
//}

// addStacktrace generates and adds stacktrace
// (if it's not presented by  ErrLetter's field).
func (e *Entry) addStacktrace() (this *Entry) {

	if e.ErrLetter == nil {
		e.LogLetter.StackTrace = ekasys.GetStackTrace(3, -1).ExcludeInternal()
	}

	return e
}
