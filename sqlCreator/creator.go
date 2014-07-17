package sqlCreator

import (
	"strings"
	"reflect"
	"log"
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

// @todo refactor and use reflext.Value instead of crap interface
func GetCreateTableByStruct(anyStruct interface{}) (string) {
	tableName, reflectValue := getTableName(anyStruct)

	if true == tableIsCreated(tableName) {
		return ""
	}

	columnDefinitions := getSqlConfigFromStruct(reflectValue, tableName)

	createTable := "CREATE TABLE " + QuoteInto(tableName)
	var columns []string
	columns = append(columns, "\t`id` varchar(50) NOT NULL")

	for i := 0; i < len(columnDefinitions.fields); i++ {
		col := columnDefinitions.fields[i]
		columns = append(columns, "\t`"+col.name+"` "+col.colType)
		tableColumnAdd(tableName, col.name)
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
