package main

// An example streaming XML parser.
// Initial Source: https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

import (
	"./onixml"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"strings"
	"time"
)

var (
	timeStart    = time.Now()
	inputFile    = flag.String("infile", "demo-availability.xml", "Input file path")
	dbHost       = flag.String("host", "127.0.0.1", "MySQL host name")
	dbDb         = flag.String("db", "test", "MySQL db name")
	dbUser       = flag.String("user", "test", "MySQL user name")
	dbPass       = flag.String("pass", "test", "MySQL password")
	tablePrefix  = flag.String("tablePrefix", "onix_", "Table name prefix")
	tableColumns = make([][]string, 50)
	dbCon        *sql.DB
	tablesInDb   = make(map[string]string)
)

func handleErr(theErr error) {
	if nil != theErr {
		panic(theErr.Error())
	}
}

func quoteInto(data string) string {
	return "`" + strings.Replace(data, "`", "", -1) + "`"
}

func getConnection() *sql.DB {
	var dbConErr error

	if nil == dbCon {
		dbCon, dbConErr = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			url.QueryEscape(*dbUser),
			url.QueryEscape(*dbPass),
			*dbHost,
			*dbDb))
		handleErr(dbConErr)
		// why is defer close not working here?
	}
	return dbCon
}

func initDatabase() {
	// Open doesn't open a connection. Validate DSN data:
	err := getConnection().Ping()
	handleErr(err)

	// delete already created tables
	// escape dbdb due to SQL injection
	columnName := quoteInto("Tables_in_" + *dbDb)
	showQuery := "SHOW TABLES FROM " + quoteInto(*dbDb) + " WHERE " + columnName + " LIKE '" + url.QueryEscape(*tablePrefix) + "%'"
	rows, err := getConnection().Query(showQuery)
	handleErr(err)
	defer rows.Close()
	for rows.Next() { // just for learning purpose otherwise we can directly drop tables here
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil { /* error handling */
		}
		tablesInDb[tableName] = tableName
	}

	if len(tablesInDb) > 0 {
		for table := range tablesInDb {
			_, err := getConnection().Query("DROP TABLE " + quoteInto(table))
			handleErr(err)
		}
		log.Printf("Dropped %d existing tables", len(tablesInDb))
	}
}

func printDuration() {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	fmt.Printf("XML Parser took %dh %dm %fs to run.\n", int(duration.Hours()), int(duration.Minutes()), duration.Seconds())
	fmt.Printf("XML Parser took %v to run.\n", duration)
}

func main() {
	flag.Parse()
	initDatabase()

	onixml.SetConnection(getConnection())
	total := onixml.OnixmlDecode(*inputFile)

	fmt.Printf("Total articles: %d \n", total)
	getConnection().Close()
	printDuration()
}
