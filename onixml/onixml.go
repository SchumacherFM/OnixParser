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
	"database/sql"
	"encoding/xml"
	"log"
	"os"
	"reflect"
	"sync" // for concurrency
	"time"
	"runtime"
)

var (
	dbCon              *sql.DB
	tablePrefix        string
	Verbose *bool
)

func SetConnection(aCon *sql.DB) {
	dbCon = aCon
}
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

func handleErr(theErr error) {
	if nil != theErr {
		panic(theErr)
		//	log.Fatal(theErr.Error())
	}
}

func logger(format string, v ...interface{}) {
	if *Verbose {
		log.Printf(format, v...)
	}
}

func OnixmlDecode(inputFile string) (int, int) {

	sqlCreator.SetTablePrefix(tablePrefix)
	total := 0
	totalErr := 0
	xmlFile, err := os.Open(inputFile)
	handleErr(err)
	xmlStat, err := xmlFile.Stat()
	handleErr(err)
	if true == xmlStat.IsDir() {
		logger("%s is a directory ...\n", inputFile)
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
		handleErr(dtErr)

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
					logger("Decode Error, Type mismatch: %v\n%v\n", prod, decErr)
					totalErr++
				}
				wg.Add(1)
				go parseXmlElementsConcurrent(&prod, dbCon, &wg)

				if total > 0 && 0 == total%1000 {
					printDuration(timeStart, total)
					timeStart = time.Now()
				}
				total++
			}
		default:
		}
	}
	c := time.Tick(10 * time.Second) // every 10 seconds
	for now := range c {
		numRoutines := runtime.NumGoroutine()
		logger("%d child processes remaining ... %v", numRoutines, now)
		if numRoutines < 10 {
			break;
		}
	}
	wg.Wait() // wait for the goroutines to finish, is that now redundant regarding the infinite for loop?

	return total, totalErr
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
	_, err := dbCon.Exec(createTable) // instead of .Query because we don't care for result. Exec closes resource
	handleErr(err)
}

func getNameOfStruct(anyStruct interface{}) string {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	return typeOfAnyStruct.Name()
}

func getInsertStmt(anyStruct interface{}) string {
	return sqlCreator.GetInsertTableByStruct(anyStruct)
}

func parseXmlElementsConcurrent(prod *Product, sharedDbCon *sql.DB, wg *sync.WaitGroup) {
	SetConnection(sharedDbCon)  // as we are in another thread set the dbCon new
	defer wg.Done()

	prod.writeToDb("")

	if len(prod.ProductIdentifier) > 0 {
		for _, prodIdentifier := range prod.ProductIdentifier {
			prodIdentifier.writeToDb(prod.RecordReference)
		}
	}
	prod.Title.writeToDb(prod.RecordReference)
	prod.Series.writeToDb(prod.RecordReference)
	prod.Website.writeToDb(prod.RecordReference)
	prod.Extent.writeToDb(prod.RecordReference)

	if len(prod.Contributor) > 0 {
		for _, prodContributor := range prod.Contributor {
			prodContributor.writeToDb(prod.RecordReference)
		}
	}

	if len(prod.Subject) > 0 {
		for _, prodSubject := range prod.Subject {
			prodSubject.writeToDb(prod.RecordReference)
		}
	}
	if len(prod.SupplyDetail) > 0 {
		for _, prodSupplyDetail := range prod.SupplyDetail {
			prodSupplyDetail.writeToDb(prod.RecordReference)
		}
	}

	// @todo convert all other methods here to struct based ones
	if len(prod.OtherText) > 0 {
		for _, prodOtherText := range prod.OtherText {
			if prodOtherText.TextTypeCode > 0 {
				xmlElementOtherText(prod.RecordReference, &prodOtherText)
			}
		}
	}

	if prod.MediaFile.MediaFileTypeCode > 0 {
		xmlElementMediaFile(prod.RecordReference, &prod.MediaFile)
	}
	if "" != prod.Imprint.ImprintName {
		xmlElementImprint(prod.RecordReference, &prod.Imprint)
	}
	if prod.Publisher.PublishingRole > 0 {
		xmlElementPublisher(prod.RecordReference, &prod.Publisher)
	}
	if prod.SalesRights.SalesRightsType > 0 {
		xmlElementSalesRights(prod.RecordReference, &prod.SalesRights)
	}
	if prod.SalesRestriction.SalesRestrictionType > 0 {
		xmlElementSalesRestriction(prod.RecordReference, &prod.SalesRestriction)
	}
	if prod.Measure.MeasureTypeCode > 0 {
		xmlElementMeasure(prod.RecordReference, &prod.Measure)
	}
	if prod.RelatedProduct.ProductIDType > 0 {
		xmlElementRelatedProduct(prod.RecordReference, &prod.RelatedProduct)
	}
	if "" != prod.MarketRepresentation.AgentName {
		xmlElementMarketRepresentation(prod.RecordReference, &prod.MarketRepresentation)
	}
}

func xmlElementOtherText(id string, o *OtherText) {
	iSql := getInsertStmt(o)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		o.TextTypeCode,
		o.Text)
	handleErr(stmtErr)
}

func xmlElementMediaFile(id string, m *MediaFile) {
	iSql := getInsertStmt(m)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		m.MediaFileTypeCode,
		m.MediaFileLinkTypeCode,
		m.MediaFileLink)
	handleErr(stmtErr)
}

func xmlElementImprint(id string, i *Imprint) {
	iSql := getInsertStmt(i)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		i.ImprintName)
	handleErr(stmtErr)
}

func xmlElementPublisher(id string, p *Publisher) {
	iSql := getInsertStmt(p)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		p.PublishingRole,
		p.PublisherName)
	handleErr(stmtErr)
}
func xmlElementSalesRights(id string, s *SalesRights) {
	iSql := getInsertStmt(s)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		s.SalesRightsType,
		s.RightsCountry)
	handleErr(stmtErr)
}
func xmlElementSalesRestriction(id string, s *SalesRestriction) {
	iSql := getInsertStmt(s)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		s.SalesRestrictionType)
	handleErr(stmtErr)
}

func xmlElementMeasure(id string, m *Measure) {
	iSql := getInsertStmt(m)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		m.MeasureTypeCode,
		m.Measurement,
		m.MeasureUnitCode)
	handleErr(stmtErr)
}

func xmlElementRelatedProduct(id string, r *RelatedProduct) {
	iSql := getInsertStmt(r)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		r.RelationCode,
		r.ProductIDType,
		r.IDValue)
	handleErr(stmtErr)
}

func xmlElementSupplyDetailPrice(id string, supplierName string, p *Price) {
	iSql := getInsertStmt(p)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		supplierName,
		p.PriceTypeCode,
		p.DiscountCodeType,
		p.DiscountCode,
		p.PriceAmount,
		p.CurrencyCode,
		p.CountryCode)
	handleErr(stmtErr)
}

func xmlElementMarketRepresentation(id string, m *MarketRepresentation) {
	iSql := getInsertStmt(m)
	_, stmtErr := dbCon.Exec(
		iSql, id,
		m.AgentName,
		m.AgentRole,
		m.MarketCountry,
		m.MarketPublishingStatus)
	handleErr(stmtErr)
}

func printDuration(timeStart time.Time, currentCount int) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	mem := float64(memStats.Alloc) / 1024 / 1024
	logger("%v Processed: %d, child processes: %d, Mem alloc: %.2fMB\n",
		duration,
		currentCount,
		runtime.NumGoroutine(),
		mem)
}
