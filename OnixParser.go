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
	"os"
	"runtime"
	"time"
)

type appConfiguration struct {
	inputFile   *string
	dbHost      *string
	dbDb        *string
	dbUser      *string
	dbPass      *string
	tablePrefix *string
	verbose     *bool
	dbCon       *sql.DB
	maxLoadAvg  *float64
	maxOpenCon  *int
}

var appConfig = appConfiguration{
	inputFile:   flag.String("infile", "", "Input file path"),
	dbHost:      flag.String("host", "127.0.0.1", "MySQL host name"),
	dbDb:        flag.String("db", "test", "MySQL db name"),
	dbUser:      flag.String("user", "test", "MySQL user name"),
	dbPass:      flag.String("pass", "test", "MySQL password"),
	tablePrefix: flag.String("tablePrefix", "gonix_", "Table name prefix"),
	verbose:     flag.Bool("v", false, "Increase verbosity"),
	maxLoadAvg:  flag.Float64("mla", 6.5, "Max Load Average, float value. Recommended > 6, if <= 3 then disabled"),
	maxOpenCon:  flag.Int("moc", 20, "Max MySQL open connections"),
}

func handleErr(theErr error) {
	if nil != theErr {
		log.Fatal(theErr.Error())
	}
}

func logger(format string, v ...interface{}) {
	if *appConfig.verbose {
		log.Printf(format, v...)
	}
}

func getConnection() *sql.DB {
	var dbConErr error

	if nil == appConfig.dbCon {
		appConfig.dbCon, dbConErr = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			url.QueryEscape(*appConfig.dbUser),
			url.QueryEscape(*appConfig.dbPass),
			*appConfig.dbHost,
			*appConfig.dbDb))
		handleErr(dbConErr)
		appConfig.dbCon.SetMaxIdleConns(5)
		appConfig.dbCon.SetMaxOpenConns(int(*appConfig.maxOpenCon)) // amount of structs
		// why is defer close not working here?
	}
	return appConfig.dbCon
}

func initDatabase() {
	// Open doesn't open a connection. Validate DSN data:
	err := getConnection().Ping()
	handleErr(err)

	// delete already created tables
	// escape dbdb due to SQL injection
	columnName := sqlCreator.QuoteInto("Tables_in_" + *appConfig.dbDb)
	showQuery := "SHOW TABLES FROM " + sqlCreator.QuoteInto(*appConfig.dbDb) + " WHERE " + columnName + " LIKE '" + url.QueryEscape(*appConfig.tablePrefix) + "%'"
	rows, err := getConnection().Query(showQuery)
	handleErr(err)
	defer rows.Close()
	rowCount := 0
	for rows.Next() { // just for learning purpose otherwise we can directly drop tables here
		var tableName string
		err = rows.Scan(&tableName)
		handleErr(err)
		_, err := getConnection().Exec("DROP TABLE " + sqlCreator.QuoteInto(tableName))
		handleErr(err)
		rowCount++
	}
	logger("Dropped %d existing tables", rowCount)
}

func printDuration(timeStart time.Time) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	logger("XML Parser took %dh %dm %fs to run.\n", int(duration.Hours()), int(duration.Minutes()), duration.Seconds())
	logger("XML Parser took %v to run.\n", duration)
}

func main() {
	timeStart := time.Now()
	if "" == os.Getenv("GOMAXPROCS") {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	fmt.Println("OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html")
	flag.Parse()
	initDatabase()
	onixml.SetAppConfig(appConfig.dbCon, appConfig.tablePrefix, appConfig.inputFile, appConfig.maxLoadAvg, appConfig.verbose)
	total, totalErr := onixml.OnixmlDecode()

	logger("Total products: %d \n", total)
	logger("Total errors: %d \n", totalErr)
	getConnection().Close()
	printDuration(timeStart)
}
