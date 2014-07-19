/*
	Copyright (C) 2014  Cyrill AT Schumacher dot fm

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

    Contribute @ https://github.com/SchumacherFM/OnixParser
*/
package sqlCreator

import (
	"reflect"
	"sync"
)

// typeInfo holds details for the xml representation of a type.
type typeInfo struct {
	structName string
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
			if f.PkgPath != "" || "" == f.Tag.Get("sql") {
				continue // Private field or sub struct
			}

			finfo, err := structFieldInfo(typ, &f)
			if err != nil {
				return nil, err
			}

			tinfo.structName = typ.Name()

			// Add the field if it doesn't conflict with other fields.
			tinfo.fields = append(tinfo.fields, *finfo)
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
