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
	"database/sql"
	"crypto/rand"
	"log"
	"net/url"
	"fmt"
)

type AppConfiguration struct {
	InputFile     *string
	OutputFile    *string
	dbHost        *string
	DbDb          *string
	dbUser        *string
	dbPass        *string
	TablePrefix   *string
	Verbose       *bool
	dbCon         *sql.DB
	MaxGoRoutines *int
	maxOpenCon    *int
}

func NewAppConfiguration() *AppConfiguration {
	a := new(AppConfiguration)
	return a
}

func (a *AppConfiguration) SetConnection(host *string, db *string, user *string, pass *string, maxOpenCon    *int) {
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

func (a *AppConfiguration) GetOutputFilePrefix() string {
	return "/tmp/" + randString(12) + ".csv"
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

func randString(n int) string {
	const alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphaNum[b%byte(len(alphaNum))]
	}
	return string(bytes)
}
