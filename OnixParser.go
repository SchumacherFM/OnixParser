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
package main

import (
	"./onixml"
	"./sqlCreator"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"time"
)

var (
	inputFile   = flag.String("infile", "demo-availability.xml", "Input file path")
	dbHost      = flag.String("host", "127.0.0.1", "MySQL host name")
	dbDb        = flag.String("db", "test", "MySQL db name")
	dbUser      = flag.String("user", "test", "MySQL user name")
	dbPass      = flag.String("pass", "test", "MySQL password")
	tablePrefix = flag.String("tablePrefix", "gonix_", "Table name prefix")
	dbCon       *sql.DB
	tablesInDb  = make(map[string]string)
)

func handleErr(theErr error) {
	if nil != theErr {
		log.Fatal(theErr.Error())
	}
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
	columnName := sqlCreator.QuoteInto("Tables_in_" + *dbDb)
	showQuery := "SHOW TABLES FROM " + sqlCreator.QuoteInto(*dbDb) + " WHERE " + columnName + " LIKE '" + url.QueryEscape(*tablePrefix) + "%'"
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
			_, err := getConnection().Exec("DROP TABLE " + sqlCreator.QuoteInto(table))
			handleErr(err)
		}
		log.Printf("Dropped %d existing tables", len(tablesInDb))
	}
}

func printDuration(timeStart time.Time) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	log.Printf("XML Parser took %dh %dm %fs to run.\n", int(duration.Hours()), int(duration.Minutes()), duration.Seconds())
	log.Printf("XML Parser took %v to run.\n", duration)
}

func main() {
	timeStart := time.Now()
	fmt.Println("OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html")
	flag.Parse()
	initDatabase()

	onixml.SetConnection(getConnection())
	onixml.SetTablePrefix(*tablePrefix)
	total, totalErr := onixml.OnixmlDecode(*inputFile)

	log.Printf("Total articles: %d \n", total)
	log.Printf("Total errors: %d \n", totalErr)
	getConnection().Close()
	printDuration(timeStart)
}
