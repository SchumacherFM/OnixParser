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
	"flag"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"time"

	_ "github.com/SchumacherFM/OnixParser/Godeps/_workspace/src/github.com/go-sql-driver/mysql"

	"github.com/SchumacherFM/OnixParser/gonfig"
	"github.com/SchumacherFM/OnixParser/onixml"
	"github.com/SchumacherFM/OnixParser/sqlCreator"
)

var appConfig = gonfig.NewAppConfiguration()

func initDatabase() {
	// Open doesn't open a connection. Validate DSN data:
	err := appConfig.GetConnection().Ping()
	appConfig.HandleErr(err)

	// delete already created tables
	// escape dbdb due to SQL injection
	columnName := sqlCreator.QuoteInto("Tables_in_" + *appConfig.DbDb)
	showQuery := "SHOW TABLES FROM " + sqlCreator.QuoteInto(*appConfig.DbDb) + " WHERE " + columnName + " LIKE '" + url.QueryEscape(*appConfig.TablePrefix) + "%'"
	rows, err := appConfig.GetConnection().Query(showQuery)
	appConfig.HandleErr(err)
	defer rows.Close()
	rowCount := 0
	for rows.Next() { // just for learning purpose otherwise we can directly drop tables here
		var tableName string
		err = rows.Scan(&tableName)
		appConfig.HandleErr(err)
		_, err := appConfig.GetConnection().Exec("DROP TABLE " + sqlCreator.QuoteInto(tableName))
		appConfig.HandleErr(err)
		rowCount++
	}
	appConfig.Log("Dropped %d existing tables", rowCount)

	// check for SELECT @@max_allowed_packet for reading LOAD DATA INFILE LOCAL
	rows2, err2 := appConfig.GetConnection().Query("SELECT @@max_allowed_packet")
	appConfig.HandleErr(err2)
	defer rows2.Close()
	if !rows2.Next() {
		rows2.Close()
		panic("crashed")
	}
	err = rows2.Scan(&appConfig.MaxPacketSize)
	appConfig.MaxPacketSize = appConfig.MaxPacketSize - 40960 // just in case remove 40kb
	appConfig.HandleErr(err)
}

func main() {
	timeStart := time.Now()
	if "" == os.Getenv("GOMAXPROCS") {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	fmt.Println("OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html")
	flag.Parse()
	appConfig.Init()
	initDatabase()
	onixml.SetAppConfig(appConfig)
	total, totalErr := onixml.OnixmlDecode()
	appConfig.Log("Total products: %d \n", total)
	appConfig.Log("Total errors: %d \n", totalErr)

	onixml.ImportCsvIntoMysql()

	appConfig.GetConnection().Close()
	appConfig.CloseOutputFiles()
	printDuration(timeStart)
}

func printDuration(timeStart time.Time) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	appConfig.Log("XML Parser took %dh %dm %fs to run.\n", int(duration.Hours()), int(duration.Minutes()), duration.Seconds())
	appConfig.Log("XML Parser took %v to run.\n", duration)
}
