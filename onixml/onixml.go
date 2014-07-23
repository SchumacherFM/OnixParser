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
	"../sqlCreator"
	"encoding/xml"
	"os"
	"reflect"
	"runtime"
	"sync" // for concurrency
	"time"
	"github.com/SchumacherFM/OnixParser/gonfig"

)

var appConfig *gonfig.AppConfiguration

// inherits from OnixParser.go
func SetAppConfig(ac *gonfig.AppConfiguration) {
	appConfig = ac
}

func initFileWriter() {

	//	outFile := "/tmp/" + *appConfig.tablePrefix + randString(10) + ".csv"
	//	if "" != *appConfig.outputFile {
	//		outFile = *appConfig.outputFile
	//	}
	//
	//	fmt.Println(outFile)

	//appConfig.fileWriter :    getFileWriterBuffer(),
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
				go parseXmlElementsConcurrent(&prod, appConfig, &wg)

				if true == *appConfig.Verbose && total > 0 && 0 == total%1000 {
					printDuration(timeStart, total)
					timeStart = time.Now()
				}
				total++
				handleAmountOfGoRoutines()
			}
		default:
		}
	}
	wg.Wait() // wait for the goroutines to finish, is that now redundant regarding the infinite for loop?
	return total, totalErr
}

func handleAmountOfGoRoutines() {
	if runtime.NumGoroutine() > *appConfig.MaxGoRoutines {
		c := time.Tick(5 * time.Second)
		for now := range c {
			ngo := runtime.NumGoroutine()
			appConfig.Log("Too many child processes: %d/%d ... %v", ngo, *appConfig.MaxGoRoutines, now)
			if ngo < *appConfig.MaxGoRoutines || ngo < 10 {
				break
			}
		}
	}
}

func printWaitForGoRoutines() {
	c := time.Tick(10 * time.Second) // every 10 seconds
	for now := range c {
		numRoutines := runtime.NumGoroutine()
		appConfig.Log("%d child processes remaining ... %v", numRoutines, now)
		if numRoutines < 10 {
			break
		}
	}
}

func createTables() {

	// is there a way to do this easier/better?
	structSlice := make([]interface{}, 19)
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
	}
}

func createTable(anyStruct interface{}) {
	createTable := sqlCreator.GetCreateTableByStruct(anyStruct)
	_, err := appConfig.GetConnection().Exec(createTable) // instead of .Query because we don't care for result. Exec closes resource
	appConfig.HandleErr(err)
}

func getNameOfStruct(anyStruct interface{}) string {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	return typeOfAnyStruct.Name()
}

func getInsertStmt(anyStruct interface{}) string {
	return sqlCreator.GetInsertTableByStruct(anyStruct)
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
		mem)
}
