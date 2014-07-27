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
	"strings"
	"sync"
)

const (
	AMOUNT_OF_STRUCTS = 19
)

var (
	outputFileNameCache   = make(map[string]string)
	outputFileNameCounter = make(map[string]int)
)

type AppConfiguration struct {
	logFile   	*string
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
	MaxPacketSize int
	Csv           struct {
		LineEnding byte
		Delimiter  byte
		Enclosure  byte
	}
}

func NewAppConfiguration() *AppConfiguration {
	a := new(AppConfiguration)
	a.logFile = flag.String("logfile", "", "Logfile name, if empty direct output")
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

	/**
	* http://en.wikipedia.org/wiki/Unit_separator#Field_separators
	* 31 Unit Separator    CSV_ENCLOSED_BY
	* 30 Record Separator  CSV_SEPARATOR
	* 29 Group Separator   EOL
	* 28 File Separator
	 */
	a.Csv.LineEnding = byte(29)
	a.Csv.Delimiter = byte(30)
	a.Csv.Enclosure = byte(31)
	a.MaxPacketSize = 1<<24-1 // 16 MB and consider this as a constant ;-)
	a.outputFiles = make(map[string]*os.File)
	return a
}

func (a *AppConfiguration) SetConnection(host *string, db *string, user *string, pass *string, maxOpenCon *int) {
	a.dbHost = host
	a.DbDb = db
	a.dbUser = user
	a.dbPass = pass
	a.maxOpenCon = maxOpenCon
}

func (a *AppConfiguration) Init(){
	if "" != *a.logFile {
		logFilePointer, err := os.OpenFile(*a.logFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFilePointer)
	}
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
	}
	return a.dbCon
}

func (a *AppConfiguration) GetOutputFileName(sqlTableName string) string {

	_, isSet := outputFileNameCache[sqlTableName]
	if true == isSet {
		return outputFileNameCache[sqlTableName]
	}

	path := *a.TablePrefix + sqlTableName + "_" + randString(12) + ".csv"
	if "" == *a.outputDir {
		outputFileNameCache[sqlTableName] = "/tmp/"+path
	} else {
		outputFileNameCache[sqlTableName] = *a.outputDir+path
	}
	return outputFileNameCache[sqlTableName]
}

func (a *AppConfiguration) InitOutputFile(tableName string) {

	fileName := a.GetOutputFileName(tableName)

	filePointer, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		a.HandleErr(err)
		os.Exit(1)
	}
	a.Lock()
	defer a.Unlock()
	a.outputFiles[tableName] = filePointer
}

func (a *AppConfiguration) GetOutputFiles() []string {
	var filesGreaterZero []string
	for tableName, filePointer := range a.outputFiles {
		fps, err := filePointer.Stat()
		a.HandleErr(err)
		if fps.Size() > 0 {
			filesGreaterZero = append(filesGreaterZero, tableName)
		}
	}
	return filesGreaterZero
}

func (a *AppConfiguration) CloseOutputFiles() {
	for _, fp := range a.outputFiles {
		err := fp.Close()
		a.HandleErr(err)
	}
}

func (a *AppConfiguration) GetOutputFilePointer(tableName string) *os.File {
	a.RLock()
	defer a.RUnlock()
	fp, ok := a.outputFiles[tableName]
	if !ok {
		panic("Failed to get file pointer for tableName: " + tableName)
	}
	return fp
}

func (a *AppConfiguration) WriteBytes(tableName string, sliceOfBytes []byte) (int, error) {
	return a.GetOutputFilePointer(tableName).Write(sliceOfBytes)
}

func (a *AppConfiguration) GetNextTableName(tableName string) string {
	_, isSet := outputFileNameCounter[tableName]
	if false == isSet {
		outputFileNameCounter[tableName] = 0
	}
	outputFileNameCounter[tableName]++
	nextTn := fmt.Sprintf("%s@%05d", tableName, outputFileNameCounter[tableName])
	a.InitOutputFile(nextTn)
	return nextTn
}

func (a *AppConfiguration) RemoveNumbersFromTableName(tableName string) string {
	return strings.Join(strings.Split(tableName, "@")[0:1], "")
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

func (a *AppConfiguration) Panic(theErr error) {
	if nil != theErr {
		panic(theErr)
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
