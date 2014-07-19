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
	"strings"
	"reflect"
	"log"
)

var (
	tablePrefix string
)

func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

func QuoteInto(data string) string {
	return "`" + strings.Replace(data, "`", "", -1) + "`"
}


func handleErr(theErr error) {
	if nil != theErr {
		log.Fatal(theErr.Error())
	}
}

func isValidTableColumn(elementType string) bool {
	if true == strings.Contains(elementType, "onixml.") {
		return false
	}
	elementTypeByte := []byte(elementType)
	elementType2 := string(elementTypeByte[0:3]) // substring(str,0,3) ;-)
	return elementType2 == "int" || elementType2 == "str" || elementType2 == "flo"
}

func getTableName(anyStruct interface{}) (string, reflect.Value) {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	return tablePrefix+strings.ToLower(typeOfAnyStruct.Name()), s
}

// @todo use Cachekey to speed up and avoid reflection
func getSqlConfigFromStruct(val reflect.Value, cacheKey string) *typeInfo {
	var tinfo        *typeInfo
	var err error
	typ := val.Type()
	tinfo, err = getTypeInfo(typ)
	handleErr(err)
	return tinfo
}

// @todo refactor and use reflext.Value instead of interface
func GetCreateTableByStruct(anyStruct interface{}) (string) {
	tableName, reflectValue := getTableName(anyStruct)

	columnDefinitions := getSqlConfigFromStruct(reflectValue, tableName)

	createTable := "CREATE TABLE " + QuoteInto(tableName)
	var columns []string
	columns = append(columns, "\t`id` varchar(50) NOT NULL")

	for i := 0; i < len(columnDefinitions.fields); i++ {
		col := columnDefinitions.fields[i]
		columns = append(columns, "\t`"+col.name+"` "+col.colType)
	}
	columns = append(columns, "\tKEY (`id`)")
	return createTable + " (\n" + strings.Join(columns, ",\n") + "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
}

func GetInsertTableByStruct(anyStruct interface{}) (string) {
	tableName, reflectValue := getTableName(anyStruct)
	columnDefinitions := getSqlConfigFromStruct(reflectValue, tableName)
	insertTable := "INSERT INTO " + QuoteInto(tableName)
	var columns, jokers []string
	columns = append(columns, "`id`")
	jokers = append(jokers, "?")
	for i := 0; i < len(columnDefinitions.fields); i++ {
		col := columnDefinitions.fields[i]
		columns = append(columns, QuoteInto(col.name))
		jokers = append(jokers, "?")
	}
	return insertTable + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(jokers, ",") + ")"
}
