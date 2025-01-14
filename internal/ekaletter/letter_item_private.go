// Copyright © 2020. All rights reserved.
// Author: Ilya Stroy.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package ekaletter

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/qioalice/ekago/v2/internal/ekaclike"
	"github.com/qioalice/ekago/v2/internal/ekafield"

	"github.com/modern-go/reflect2"
)

var (
	reflectedTimeTime = reflect2.TypeOf(time.Time{})
	reflectedTimeDuration = reflect2.TypeOf(time.Duration(0))
)

// addImplicitField adds new field to l.Fields treating 'name' as field's name,
// 'value' as field's value and using 'typ' (assuming it's value's type)
// to recognize how to convert Golang's interface{} to the 'Field' object.
func (li *LetterItem) addImplicitField(name string, value interface{}, typ reflect2.Type) {

	varyField := name != "" && name[len(name)-1] == '?'
	if varyField {
		name = name[:len(name)-1]
	}

	var f ekafield.Field

	switch {
	case value == nil && varyField:
		// do nothing
		return

	case value == nil:
		li.Fields = append(li.Fields, ekafield.NilValue(name, ekafield.KIND_TYPE_INVALID))
		return

	case typ == reflectedTimeTime:
		var timeVal time.Time
		typ.UnsafeSet(unsafe.Pointer(&timeVal), reflect2.PtrOf(value))
		f = ekafield.Time(name, timeVal)
		goto recognizer

	case typ == reflectedTimeDuration:
		var durationVal time.Duration
		typ.UnsafeSet(unsafe.Pointer(&durationVal), reflect2.PtrOf(value))
		f = ekafield.Duration(name, durationVal)
		goto recognizer

	// PLACE TYPES ABOVE THAT HAS String() METHOD BUT YOU DON'T WANT TO USE IT.

	case typ.Implements(ekafield.ReflectedTypeFmtStringer):
		f = ekafield.Stringer(name, value.(fmt.Stringer))
		goto recognizer
	}

	switch typ.Kind() {

	case reflect.Ptr:
		// Maybe it's typed nil pointer? We can't dereference nil pointer
		// but I guess it's important to log nil pointer as is even if
		// FLAG_ALLOW_IMPLICIT_POINTERS is not set (because what can we do otherwise?)
		logPtrAsIs :=
			li.Flags.TestAll(FLAG_ALLOW_IMPLICIT_POINTERS) ||
				ekaclike.TakeRealAddr(value) == nil

		if logPtrAsIs {
			f = ekafield.Addr(name, value)
		} else {
			value = typ.Indirect(value)
			li.addImplicitField(name, value, reflect2.TypeOf(value))
		}

	case reflect.Bool:
		var boolVal bool
		typ.UnsafeSet(unsafe.Pointer(&boolVal), reflect2.PtrOf(value))
		f = ekafield.Bool(name, boolVal)

	case reflect.Int:
		var intVal int
		typ.UnsafeSet(unsafe.Pointer(&intVal), reflect2.PtrOf(value))
		f = ekafield.Int(name, intVal)

	case reflect.Int8:
		var int8Val int8
		typ.UnsafeSet(unsafe.Pointer(&int8Val), reflect2.PtrOf(value))
		f = ekafield.Int8(name, int8Val)

	case reflect.Int16:
		var int16Val int16
		typ.UnsafeSet(unsafe.Pointer(&int16Val), reflect2.PtrOf(value))
		f = ekafield.Int16(name, int16Val)

	case reflect.Int32:
		var int32Val int32
		typ.UnsafeSet(unsafe.Pointer(&int32Val), reflect2.PtrOf(value))
		f = ekafield.Int32(name, int32Val)

	case reflect.Int64:
		var int64Val int64
		typ.UnsafeSet(unsafe.Pointer(&int64Val), reflect2.PtrOf(value))
		f = ekafield.Int64(name, int64Val)

	case reflect.Uint:
		var uintVal uint64
		typ.UnsafeSet(unsafe.Pointer(&uintVal), reflect2.PtrOf(value))
		f = ekafield.Uint64(name, uintVal)

	case reflect.Uint8:
		var uint8Val uint8
		typ.UnsafeSet(unsafe.Pointer(&uint8Val), reflect2.PtrOf(value))
		f = ekafield.Uint8(name, uint8Val)

	case reflect.Uint16:
		var uint16Val uint16
		typ.UnsafeSet(unsafe.Pointer(&uint16Val), reflect2.PtrOf(value))
		f = ekafield.Uint16(name, uint16Val)

	case reflect.Uint32:
		var uint32Val uint32
		typ.UnsafeSet(unsafe.Pointer(&uint32Val), reflect2.PtrOf(value))
		f = ekafield.Uint32(name, uint32Val)

	case reflect.Uint64:
		var uint64Val uint64
		typ.UnsafeSet(unsafe.Pointer(&uint64Val), reflect2.PtrOf(value))
		f = ekafield.Uint64(name, uint64Val)

	case reflect.Float32:
		var float32Val float32
		typ.UnsafeSet(unsafe.Pointer(&float32Val), reflect2.PtrOf(value))
		f = ekafield.Float32(name, float32Val)

	case reflect.Float64:
		var float64Val float64
		typ.UnsafeSet(unsafe.Pointer(&float64Val), reflect2.PtrOf(value))
		f = ekafield.Float64(name, float64Val)

	case reflect.Complex64:
		var complex64Val complex64
		typ.UnsafeSet(unsafe.Pointer(&complex64Val), reflect2.PtrOf(value))
		f = ekafield.Complex64(name, complex64Val)

	case reflect.Complex128:
		var complex128Val complex128
		typ.UnsafeSet(unsafe.Pointer(&complex128Val), reflect2.PtrOf(value))
		f = ekafield.Complex128(name, complex128Val)

	case reflect.String:
		var stringVal string
		typ.UnsafeSet(unsafe.Pointer(&stringVal), reflect2.PtrOf(value))
		f = ekafield.String(name, stringVal)

	case reflect.Uintptr, reflect.UnsafePointer:
		f = ekafield.Addr(name, value)

	// TODO: handle all structs, handle structs with Valid (bool) = false as null

	default:
	}

recognizer:
	if !(varyField && f.IsZero()) {
		li.Fields = append(li.Fields, f)
	}
}

// addExplicitFieldByPtr adds 'f' to the l.Fields only if it's not nil and
// if it's not a vary-zero field.
func (li *LetterItem) addExplicitFieldByPtr(f *ekafield.Field) {
	if f != nil {
		li.addExplicitField2(*f)
	}
}

// addExplicitField2 adds 'f' to the l.Fields only if it's not vary-zero field.
func (li *LetterItem) addExplicitField2(f ekafield.Field) {
	varyField := f.Key != "" && f.Key[len(f.Key)-1] == '?'
	if varyField {
		f.Key = f.Key[:len(f.Key)-1]
	}
	if !(varyField && f.IsZero()) {
		li.Fields = append(li.Fields, f)
	}
}
