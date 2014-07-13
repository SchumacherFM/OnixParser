package sqlCreator

import (
	"fmt"
	"strings"
	"reflect"
)

type tableColumn struct {
	Table, Column string
}

var (
	tablePrefix string
	tableColumns = make(map[tableColumn]bool) // key table,column name = value bool
)

func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

func QuoteInto(data string) string {
	return "`" + strings.Replace(data, "`", "", -1) + "`"
}


func handleErr(theErr error) {
	if nil != theErr {
		panic(theErr.Error())
	}
}

func isValidTableColumn(elementType string) bool {
	if true == strings.Contains(elementType, "onixml.") {
		return false
	}
	elementTypeByte := []byte(elementType)
	elementType2 := string(elementTypeByte[0:3])
	return elementType2 == "int" || elementType2 == "str" || elementType2 == "flo"
}


func getTableColumn(elementType string) string {

	if false == isValidTableColumn(elementType) {
		return ""
	}

	elementTypeByte := []byte(elementType)
	elementType2 := string(elementTypeByte[0:3])
	switch elementType2 {
	case "int":
		return "bigint(14) NOT NULL DEFAULT 0"
	case "str":
		return "text NULL"
	case "flo":
		return "decimal(10,2) NOT NULL DEFAULT 0.0"
	default:
		//fmt.Printf("type %s => %s not supported\n", elementType2, elementType)
	}
	return ""
}

func getTableName(name string) string {
	return tablePrefix + strings.ToLower(name)
}

func GetCreateTableByStruct(anyStruct interface{}) (string) {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	tableName := getTableName(typeOfAnyStruct.Name())

	if true == tableIsCreated(tableName) {
		return ""
	}

	createTable := "CREATE TABLE " + QuoteInto(tableName)
	var columns []string
	columns = append(columns, "\t`id` varchar(50) NOT NULL")

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		//				fmt.Printf("%d: %s %s = %v\n", i,
		//					typeOfAnyStruct.Field(i).Name, f.Type(), f.Interface())
		columnName := typeOfAnyStruct.Field(i).Name
		columnSqlType := getTableColumn(f.Type().Name())
		if "" != columnSqlType {
			columns = append(columns, fmt.Sprintf("\t`%s` %s", columnName, columnSqlType))
			tableColumnAdd(tableName, columnName)
		}
	}
	columns = append(columns, "\tKEY (`id`)")
	return createTable + " (\n" + strings.Join(columns, ",\n") + "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
}

func tableColumnAdd(tableName string, columnName string) {
	tableColumns[tableColumn{tableName, ""}] = true
	tableColumns[tableColumn{tableName, columnName}] = true
}

func tableIsCreated(tableName string) bool {
	t := tableColumn{tableName, ""}
	_, isSet := tableColumns[t]
	return isSet
}

func GetInsertTableByStruct(anyStruct interface{}) (string) {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	tableName := getTableName(typeOfAnyStruct.Name())

	insertTable := "INSERT INTO " + QuoteInto(tableName)
	var columns, jokers []string
	columns = append(columns, "`id`")
	jokers = append(jokers, "?")
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		columnName := typeOfAnyStruct.Field(i).Name

		if false == isValidTableColumn(f.Type().Name()) { // if type of column
			continue
		}
		columns = append(columns, QuoteInto(columnName))
		jokers = append(jokers, "?")
	}
	// @todo check for column miss match
	return insertTable + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(jokers, ",") + ")"
}
