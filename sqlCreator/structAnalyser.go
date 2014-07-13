// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shamelessly copied from the encoding/xml package
// modified to read the sql part from the backtick section

package sqlCreator

import (
	//	"fmt"
	//	"strings"
	"reflect"
	"sync"
)

// typeInfo holds details for the xml representation of a type.
type typeInfo struct {
	xmlname string
	fields  []fieldInfo
}

// fieldInfo holds details for the xml representation of a single field.
type fieldInfo struct {
	idx       []int
	name      string
	colType   string
}

type Name struct {
	Space, Local string
}

type fieldFlags int

var tinfoMap = make(map[reflect.Type]*typeInfo)
var tinfoLock sync.RWMutex

var nameType = reflect.TypeOf(Name{})

// getTypeInfo returns the typeInfo structure with details necessary
// for marshalling and unmarshalling typ.
func getTypeInfo(typ reflect.Type) (*typeInfo, error) {
	tinfoLock.RLock()
	tinfo, ok := tinfoMap[typ]
	tinfoLock.RUnlock()
	if ok {
		return tinfo, nil
	}
	tinfo = &typeInfo{}
	if typ.Kind() == reflect.Struct && typ != nameType {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.PkgPath != "" || f.Tag.Get("sql") == "-" {
				continue // Private field
			}

			finfo, err := structFieldInfo(typ, &f)
			if err != nil {
				return nil, err
			}

			tinfo.xmlname = typ.Name()

			// Add the field if it doesn't conflict with other fields.
			if err := addFieldInfo(typ, tinfo, finfo); err != nil {
				return nil, err
			}
		}
	}
	tinfoLock.Lock()
	tinfoMap[typ] = tinfo
	tinfoLock.Unlock()
	return tinfo, nil
}

// structFieldInfo builds and returns a fieldInfo for f.
func structFieldInfo(typ reflect.Type, f *reflect.StructField) (*fieldInfo, error) {
	finfo := &fieldInfo{idx: f.Index}

	tag := f.Tag.Get("sql")
	finfo.colType = tag
	finfo.name = f.Name

	return finfo, nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// addFieldInfo adds finfo to tinfo.fields if there are no
// conflicts, or if conflicts arise from previous fields that were
// obtained from deeper embedded structures than finfo. In the latter
// case, the conflicting entries are dropped.
// A conflict occurs when the path (parent + name) to a field is
// itself a prefix of another path, or when two paths match exactly.
// It is okay for field paths to share a common, shorter prefix.
func addFieldInfo(typ reflect.Type, tinfo *typeInfo, newf *fieldInfo) error {

	tinfo.fields = append(tinfo.fields, *newf)
	return nil

}

// value returns v's field value corresponding to finfo.
// It's equivalent to v.FieldByIndex(finfo.idx), but initializes
// and dereferences pointers as necessary.
func (finfo *fieldInfo) value(v reflect.Value) reflect.Value {
	for i, x := range finfo.idx {
		if i > 0 {
			t := v.Type()
			if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}
