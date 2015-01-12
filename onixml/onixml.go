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
package onixml

import (
	"encoding/xml"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/SchumacherFM/OnixParser/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/SchumacherFM/OnixParser/gonfig"
	. "github.com/SchumacherFM/OnixParser/onixStructs"
	"github.com/SchumacherFM/OnixParser/sqlCreator"
)

var appConfig *gonfig.AppConfiguration

// inherits from OnixParser.go
func SetAppConfig(ac *gonfig.AppConfiguration) {
	appConfig = ac
}

func OnixmlDecode() (int, int) {
	sqlCreator.SetTablePrefix(appConfig.TablePrefix)
	total := 0
	totalErr := 0

	if "" == *appConfig.InputFile {
		appConfig.Log("Input file is empty\n")
		return -1, -1
	}

	xmlFile, err := os.Open(*appConfig.InputFile)
	appConfig.HandleErr(err)
	xmlStat, err := xmlFile.Stat()
	appConfig.HandleErr(err)
	if true == xmlStat.IsDir() {
		appConfig.Log("%s is a directory ...\n", appConfig.InputFile)
		return -1, -1
	}

	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	createTables()

	var wg sync.WaitGroup
	var inElement string
	timeStart := time.Now()
	for {
		// Read tokens from the XML document in a stream.
		t, dtErr := decoder.Token()
		if t == nil {
			break
		}
		appConfig.HandleErr(dtErr)

		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "Product"
			if inElement == "Product" {
				var prod Product
				// decode a whole chunk of following XML into the
				// variable prod which is a Product (se above)
				decErr := decoder.DecodeElement(&prod, &se)
				if nil != decErr {
					appConfig.Log("Decode Error, Type mismatch: %v\n%v\n", prod, decErr)
					totalErr++
				}
				wg.Add(1)
				// go here does not really make sense ... but for learning it is ok
				go ParseXmlElementsConcurrent(&prod, appConfig, &wg)

				if true == *appConfig.Verbose && total > 0 && 0 == total%1000 {
					printDuration(timeStart, total)
					timeStart = time.Now()
				}
				total++
			}
		default:
		}
	}
	wg.Wait() // wait for the goroutines to finish, is that now redundant regarding the infinite for loop?
	return total, totalErr
}

func createTables() {

	// is there a way to do this easier/better?
	structSlice := make([]interface{}, gonfig.AMOUNT_OF_STRUCTS)
	structSlice[0] = new(Product)
	structSlice[1] = new(ProductIdentifier)
	structSlice[2] = new(Title)
	structSlice[3] = new(Series)
	structSlice[4] = new(Website)
	structSlice[5] = new(Contributor)
	structSlice[6] = new(Subject)
	structSlice[7] = new(Extent)
	structSlice[8] = new(OtherText)
	structSlice[9] = new(MediaFile)
	structSlice[10] = new(Imprint)
	structSlice[11] = new(Publisher)
	structSlice[12] = new(SalesRights)
	structSlice[13] = new(SalesRestriction)
	structSlice[14] = new(Measure)
	structSlice[15] = new(RelatedProduct)
	structSlice[16] = new(SupplyDetail)
	structSlice[17] = new(Price)
	structSlice[18] = new(MarketRepresentation)

	for _, theStruct := range structSlice {
		createTable(theStruct)
		initFileWriter(theStruct)
	}
}

func initFileWriter(anyStruct interface{}) {
	tableName := appConfig.GetNameOfStruct(anyStruct)
	appConfig.InitOutputFile(tableName)
}

func createTable(anyStruct interface{}) {
	createTable := sqlCreator.GetCreateTableByStruct(anyStruct)
	_, err := appConfig.GetConnection().Exec(createTable) // instead of .Query because we don't care for result. Exec closes resource
	appConfig.HandleErr(err)
}

func ImportCsvIntoMysql() {
	infileTpl := "LOAD DATA LOCAL INFILE '%s' INTO TABLE `%s` FIELDS TERMINATED BY '%c' ENCLOSED BY '%c' LINES TERMINATED BY '%c'"

	for _, tableName := range appConfig.GetOutputFiles() {
		fileName := appConfig.GetOutputFileName(tableName)
		mysql.RegisterLocalFile(fileName)

		infileStmt := fmt.Sprintf(
			infileTpl,
			fileName,
			sqlCreator.GetTableName(appConfig.RemoveNumbersFromTableName(tableName)),
			appConfig.Csv.Delimiter,
			appConfig.Csv.Enclosure,
			appConfig.Csv.LineEnding,
		)

		appConfig.Log(infileStmt)
		_, dbErr := appConfig.GetConnection().Exec(infileStmt)
		appConfig.Panic(dbErr)
	}
}

func printDuration(timeStart time.Time, currentCount int) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	mem := float64(memStats.Sys) / 1024 / 1024
	appConfig.Log("%v Processed: %d, child processes: %d, Mem alloc: %.2fMB\n",
		duration,
		currentCount,
		runtime.NumGoroutine(),
		mem,
	)
}
