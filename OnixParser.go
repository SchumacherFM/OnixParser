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

	_ "github.com/go-sql-driver/mysql"

	"github.com/SchumacherFM/OnixParser/gonfig"
	"./onixml"
	"./sqlCreator"

)

var appConfig = gonfig.NewAppConfiguration()

func init() {
	appConfig.InputFile = flag.String("infile", "", "Input file path")
	appConfig.OutputFile = flag.String("outfile", "", "Prefix of CSV output file for reading into MySQL, if empty writes to /tmp/rand_[table].csv")
	//	appConfig.DbHost = flag.String("host", "127.0.0.1", "MySQL host name")
	//	appConfig.DbDb = flag.String("db", "test", "MySQL db name")
	//	appConfig.DbUser = flag.String("user", "test", "MySQL user name")
	//	appConfig.DbPass = flag.String("pass", "test", "MySQL password")

	appConfig.SetConnection(
		flag.String("host", "127.0.0.1", "MySQL host name"),
		flag.String("db", "test", "MySQL db name"),
		flag.String("user", "test", "MySQL user name"),
		flag.String("pass", "test", "MySQL password"),
		flag.Int("moc", 20, "Max MySQL open connections"),
	)

	appConfig.TablePrefix = flag.String("tablePrefix", "gonix_", "Table name prefix")
	appConfig.Verbose = flag.Bool("v", false, "Increase verbosity")
	appConfig.MaxGoRoutines = flag.Int("children", 2200, "Max number of sub processes. This can be up to the amount of products your importing.")

}

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
	onixml.SetAppConfig(appConfig)
	total, totalErr := onixml.OnixmlDecode()

	appConfig.Log("Total products: %d \n", total)
	appConfig.Log("Total errors: %d \n", totalErr)
	appConfig.GetConnection().Close()
	printDuration(timeStart)
}


func printDuration(timeStart time.Time) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	appConfig.Log("XML Parser took %dh %dm %fs to run.\n", int(duration.Hours()), int(duration.Minutes()), duration.Seconds())
	appConfig.Log("XML Parser took %v to run.\n", duration)
}
