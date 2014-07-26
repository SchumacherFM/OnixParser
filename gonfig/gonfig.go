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
package gonfig

import (
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"sync"
)

const (
	AMOUNT_OF_STRUCTS = 19
)

type AppConfiguration struct {
	InputFile   *string
	outputDir   *string
	outputFiles map[string]*os.File
	sync.RWMutex
	dbHost        *string
	DbDb          *string
	dbUser        *string
	dbPass        *string
	TablePrefix   *string
	Verbose       *bool
	dbCon         *sql.DB
	maxOpenCon    *int
	Csv           struct {
		LineEnding *string
		Delimiter  *string
		Enclosure  *string
		Escape     *string
	}
}

func NewAppConfiguration() *AppConfiguration {
	a := new(AppConfiguration)
	a.InputFile = flag.String("infile", "", "Input file path")
	a.outputDir = flag.String("outdir", "", "Dir for CSV output file for reading into MySQL, if empty writes to /tmp/")
	a.SetConnection(
		flag.String("host", "127.0.0.1", "MySQL host name"),
		flag.String("db", "test", "MySQL db name"),
		flag.String("user", "test", "MySQL user name"),
		flag.String("pass", "test", "MySQL password"),
		flag.Int("moc", 20, "Max MySQL open connections"),
	)
	a.TablePrefix = flag.String("tablePrefix", "gonix_", "Table name prefix")
	a.Verbose = flag.Bool("v", false, "Increase verbosity")

	a.Csv.LineEnding = flag.String("csv-l", "\n", "CSV Line Ending")
	a.Csv.Delimiter = flag.String("csv-d", "|", "CSV field delimiter")
	a.Csv.Enclosure = flag.String("csv-en", "\"", "CSV Set the field enclosure character (one character only).")
	a.Csv.Escape = flag.String("csv-es", "\\", "CSV Set the escape character (one character only). Defaults as a backslash.")

	a.outputFiles = make(map[string]*os.File, AMOUNT_OF_STRUCTS)
	return a
}

func (a *AppConfiguration) SetConnection(host *string, db *string, user *string, pass *string, maxOpenCon *int) {
	a.dbHost = host
	a.DbDb = db
	a.dbUser = user
	a.dbPass = pass
	a.maxOpenCon = maxOpenCon
}

func (a *AppConfiguration) GetConnection() *sql.DB {
	var dbConErr error

	if nil == a.dbCon {
		a.dbCon, dbConErr = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
			url.QueryEscape(*a.dbUser),
			url.QueryEscape(*a.dbPass),
			*a.dbHost,
			*a.DbDb))
		a.HandleErr(dbConErr)
		a.dbCon.SetMaxIdleConns(5)
		a.dbCon.SetMaxOpenConns(int(*a.maxOpenCon)) // amount of structs
		// why is defer close not working here?
	}
	return a.dbCon
}

func (a *AppConfiguration) getOutputFileName(sqlTableName string) string {
	path := *a.TablePrefix + sqlTableName + "_" + randString(12) + ".csv"
	if "" == *a.outputDir {
		return "/tmp/" + path
	}
	return *a.outputDir + path
}

func (a *AppConfiguration) InitOutputFile(tableName string) {

	fileName := a.getOutputFileName(tableName)

	filePointer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		a.HandleErr(err)
		os.Exit(1)
	}
	a.Lock()
	defer a.Unlock()
	a.outputFiles[tableName] = filePointer
}

func (a *AppConfiguration) CloseOutputFiles() {
	for tn, fp := range a.outputFiles {
		fp.Close()
		a.Log("Closed file: %s\n", tn)
	}
}

func (a *AppConfiguration) getOutputFilePointer(tableName string) *os.File {
	a.RLock()
	defer a.RUnlock()
	fp, ok := a.outputFiles[tableName]
	if !ok {
		panic("Failed to get file pointer for tableName: " + tableName)
	}
	return fp
}

func (a *AppConfiguration) WriteBytes(tableName string, byteString []byte) (int, error) {
	return a.getOutputFilePointer(tableName).Write(byteString)
}

func (a *AppConfiguration) Log(format string, v ...interface{}) {
	if *a.Verbose {
		log.Printf(format, v...)
	}
}

func (a *AppConfiguration) HandleErr(theErr error) {
	if nil != theErr {
		log.Fatal(theErr.Error())
	}
}

func (a *AppConfiguration) GetNameOfStruct(anyStruct interface{}) string {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	return typeOfAnyStruct.Name()
}

func randString(n int) string {
	const alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}
	return string(bytes)
}

