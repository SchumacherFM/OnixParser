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
	"io/ioutil"
	"os"
	"testing"

	"github.com/SchumacherFM/OnixParser/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

var appConfig = NewAppConfiguration()

func TestNewAppConfiguration(t *testing.T) {

	testAC := &AppConfiguration{}
	assert.IsType(t, testAC, appConfig, "Ups expected AppConfiguration object")
}

func TestCSV(t *testing.T) {

	assert.Equal(t, appConfig.Csv.LineEnding, byte(29), "LineEnding Should be equal")
	assert.Equal(t, appConfig.Csv.Delimiter, byte(30), "Delimiter Should be equal")
	assert.Equal(t, appConfig.Csv.Enclosure, byte(31), "Enclosure Should be equal")
}

func TestSetConnection(t *testing.T) {
	host := "host"
	db := "db"
	user := "user"
	pass := "pass"
	maxCon := 4711
	appConfig.SetConnection(&host, &db, &user, &pass, &maxCon)
	assert.Exactly(t, appConfig.dbHost, &host, "Host should be equal")
	assert.Exactly(t, appConfig.dbUser, &user, "User should be equal")
	assert.Exactly(t, appConfig.dbPass, &pass, "Pass should be equal")
	assert.Exactly(t, appConfig.maxOpenCon, &maxCon, "Pass should be equal")
}

func TestLoggingToFile(t *testing.T) {
	logfileName := os.TempDir() + "onixparser_test_TestLoggingToFile.log"
	appConfig.logFile = &logfileName
	isTrue := true
	appConfig.Verbose = &isTrue
	appConfig.Init()
	appConfig.Log("TestLoggingToFile")

	logFileData, err1 := ioutil.ReadFile(logfileName)
	assert.Nil(t, err1, "Cannot open/read file: "+logfileName)

	assert.Contains(t, string(logFileData), "TestLoggingToFile", "Logfile "+logfileName+" should contain the string TestLoggingToFile")
	err2 := os.Remove(logfileName)
	assert.Nil(t, err2, "Cannot remove file: "+logfileName)

}
