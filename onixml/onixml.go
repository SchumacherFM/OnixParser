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
	"time"
	"database/sql"
	"os"
	"reflect"
	"strings"
	"../sqlCreator"
	"log"
	"runtime"
)

type tableColumn struct {
	Table, Column string
}

var (
	dbCon *sql.DB
	tablePrefix string
	preparedInsertStmt = make(map[string]*sql.Stmt)
	Verbose *bool
)

func logger(format string, v ...interface{}) {
	if *Verbose {
		log.Printf(format, v...)
	}
}
func SetConnection(aCon *sql.DB) {
	dbCon = aCon
}
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

func handleErr(theErr error) {
	if nil != theErr {
		log.Fatal(theErr.Error())
	}
}

func printDuration(timeStart time.Time, objectCount int, currentCount int) {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	mem := float64(memStats.Sys) / 1024 / 1024
	logger("%v for %d entities. Processed %d, Mem alloc: %.2fMB\n", duration, objectCount, currentCount, mem)
}


func OnixmlDecode(inputFile string) (int, int) {
	sqlCreator.SetTablePrefix(tablePrefix)
	total := 0
	totalErr := 0
	xmlFile, err := os.Open(inputFile)
	handleErr(err)
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	var inElement string
	timeStart := time.Now()
	//	productChan := make(chan *Product)
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
					logger("Decode Error, Type mismatch: %v\n", prod)
					totalErr++
				}

				//				productChan <- &prod
				//				go parseXmlElements(productChan)
				// go printSomething(&prod)
				parseXmlElements(&prod)

				if total > 0 && 0 == total%1000 {
					printDuration(timeStart, 1000, total)
					timeStart = time.Now()
				}
				total++
			}
		default:
		}
	}

	// close statements
	for _, stmt := range preparedInsertStmt {
		err := stmt.Close()
		handleErr(err)
	}

	return total, totalErr
}

func parseXmlElements(prod *Product) {

	xmlElementProduct(prod)

	if len(prod.ProductIdentifier) > 0 {
		for _, prodIdentifier := range prod.ProductIdentifier {
			xmlElementProductIdentifier(prod.RecordReference, &prodIdentifier)
		}
	}
	if prod.Title.TitleType > 0 {
		xmlElementTitle(prod.RecordReference, &prod.Title)
	}
	if "" != prod.Series.TitleOfSeries || "" != prod.Series.NumberWithinSeries {
		xmlElementSeries(prod.RecordReference, &prod.Series)
	}
	if "" != prod.Website.WebsiteLink {
		xmlElementWebsite(prod.RecordReference, &prod.Website)
	}
	if len(prod.Contributor) > 0 {
		for _, prodContributor := range prod.Contributor {
			if prodContributor.SequenceNumber > 0 {
				xmlElementContributor(prod.RecordReference, &prodContributor)
			}
		}
	}
	if len(prod.Subject) > 0 {
		for _, prodSubject := range prod.Subject {
			if prodSubject.SubjectSchemeIdentifier > 0 {
				xmlElementSubject(prod.RecordReference, &prodSubject)
			}
		}
	}
	if prod.Extent.ExtentType > 0 {
		xmlElementExtent(prod.RecordReference, &prod.Extent)
	}

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

	if len(prod.SupplyDetail) > 0 {
		for _, prodSupplyDetail := range prod.SupplyDetail {
			if "" != prodSupplyDetail.SupplierName {
				xmlElementSupplyDetail(prod.RecordReference, &prodSupplyDetail)
			}
		}
	}
	if "" != prod.MarketRepresentation.AgentName {
		xmlElementMarketRepresentation(prod.RecordReference, &prod.MarketRepresentation)
	}
}

func getNameOfStruct(anyStruct interface{}) string {
	s := reflect.ValueOf(anyStruct).Elem()
	typeOfAnyStruct := s.Type()
	return typeOfAnyStruct.Name()
}

func createTable(anyStruct interface{}) {
	createTable := sqlCreator.GetCreateTableByStruct(anyStruct)
	if "" != createTable {
		_, err := dbCon.Exec(createTable) // instead of .Query because we dont care for result. Exec closes resource
		handleErr(err)
	}
}

func getInsertStmt(anyStruct interface{}) *sql.Stmt {
	structName := getNameOfStruct(anyStruct)
	_, isSet := preparedInsertStmt[structName]
	if false == isSet {
		var err error
		insertTable := sqlCreator.GetInsertTableByStruct(anyStruct)
		preparedInsertStmt[structName], err = dbCon.Prepare(insertTable)
		handleErr(err)
	}
	return preparedInsertStmt[structName]
}

func xmlElementProduct(prod *Product) {
	createTable(prod)
	insertStmt := getInsertStmt(prod)
	// _, stmtErr := insertStmt.Exec.Call(prod) => that would be nice ... but how?
	// static typed language and that would cost performance
	/* sometimes number can be 1,234 */
	_, stmtErr := insertStmt.Exec(
		// avoiding reflection
		prod.RecordReference,
		prod.RecordReference,
		prod.NotificationType,
		prod.ProductForm,
		prod.ProductFormDetail,
		prod.EditionNumber,
		strings.Replace(prod.NumberOfPages, ",", "", -1),
		prod.IllustrationsNote,
		prod.BICMainSubject,
		prod.AudienceCode,
		prod.PublishingStatus,
		prod.PublicationDate,
		prod.YearFirstPublished)
	handleErr(stmtErr)

}

func xmlElementProductIdentifier(id string, p *ProductIdentifier) {
	createTable(p)
	insertStmt := getInsertStmt(p)

	_, stmtErr := insertStmt.Exec(
		id,
		p.ProductIDType,
		p.IDValue)
	handleErr(stmtErr)
}

func xmlElementTitle(id string, t *Title) {
	createTable(t)
	insertStmt := getInsertStmt(t)
	_, stmtErr := insertStmt.Exec(
		id,
		t.TitleType,
		t.TitleText,
		t.TitlePrefix,
		t.TitleWithoutPrefix)
	handleErr(stmtErr)
}

func xmlElementSeries(id string, s *Series) {
	createTable(s)
	insertStmt := getInsertStmt(s)
	_, stmtErr := insertStmt.Exec(
		id,
		s.TitleOfSeries,
		s.NumberWithinSeries)
	handleErr(stmtErr)
}

func xmlElementWebsite(id string, w *Website) {
	createTable(w)
	insertStmt := getInsertStmt(w)
	_, stmtErr := insertStmt.Exec(
		id,
		w.WebsiteLink)
	handleErr(stmtErr)
}

func xmlElementContributor(id string, c *Contributor) {
	createTable(c)
	insertStmt := getInsertStmt(c)
	_, stmtErr := insertStmt.Exec(
		id,
		c.SequenceNumber,
		c.ContributorRole,
		c.PersonNameInverted,
		c.TitlesBeforeNames,
		c.KeyNames)
	handleErr(stmtErr)
}

func xmlElementSubject(id string, s *Subject) {
	createTable(s)
	insertStmt := getInsertStmt(s)
	_, stmtErr := insertStmt.Exec(
		id,
		s.SubjectSchemeIdentifier,
		s.SubjectCode)
	handleErr(stmtErr)
}

func xmlElementExtent(id string, e *Extent) {
	createTable(e)
	insertStmt := getInsertStmt(e)
	_, stmtErr := insertStmt.Exec(
		id,
		e.ExtentType,
		e.ExtentValue,
		e.ExtentUnit)
	handleErr(stmtErr)
}

func xmlElementOtherText(id string, o *OtherText) {
	createTable(o)
	insertStmt := getInsertStmt(o)
	_, stmtErr := insertStmt.Exec(
		id,
		o.TextTypeCode,
		o.Text)
	handleErr(stmtErr)
}

func xmlElementMediaFile(id string, m *MediaFile) {
	createTable(m)
	insertStmt := getInsertStmt(m)
	_, stmtErr := insertStmt.Exec(
		id,
		m.MediaFileTypeCode,
		m.MediaFileLinkTypeCode,
		m.MediaFileLink)
	handleErr(stmtErr)
}

func xmlElementImprint(id string, i *Imprint) {
	createTable(i)
	insertStmt := getInsertStmt(i)
	_, stmtErr := insertStmt.Exec(
		id,
		i.ImprintName)
	handleErr(stmtErr)
}

func xmlElementPublisher(id string, p *Publisher) {
	createTable(p)
	insertStmt := getInsertStmt(p)
	_, stmtErr := insertStmt.Exec(
		id,
		p.PublishingRole,
		p.PublisherName)
	handleErr(stmtErr)
}
func xmlElementSalesRights(id string, s *SalesRights) {
	createTable(s)
	insertStmt := getInsertStmt(s)
	_, stmtErr := insertStmt.Exec(
		id,
		s.SalesRightsType,
		s.RightsCountry)
	handleErr(stmtErr)
}
func xmlElementSalesRestriction(id string, s *SalesRestriction) {
	createTable(s)
	insertStmt := getInsertStmt(s)
	_, stmtErr := insertStmt.Exec(
		id,
		s.SalesRestrictionType)
	handleErr(stmtErr)
}

func xmlElementMeasure(id string, m *Measure) {
	createTable(m)
	insertStmt := getInsertStmt(m)
	_, stmtErr := insertStmt.Exec(
		id,
		m.MeasureTypeCode,
		m.Measurement,
		m.MeasureUnitCode)
	handleErr(stmtErr)
}

func xmlElementRelatedProduct(id string, r *RelatedProduct) {
	createTable(r)
	insertStmt := getInsertStmt(r)
	_, stmtErr := insertStmt.Exec(
		id,
		r.RelationCode,
		r.ProductIDType,
		r.IDValue)
	handleErr(stmtErr)
}


func xmlElementSupplyDetail(id string, s *SupplyDetail) {
	createTable(s)
	insertStmt := getInsertStmt(s)
	_, stmtErr := insertStmt.Exec(
		id,
		s.SupplierName,
		s.SupplierRole,
		s.SupplyToCountry,
		s.ProductAvailability,
		s.ExpectedShipDate,
		s.OnHand,
		s.OnOrder,
		s.PackQuantity)
	handleErr(stmtErr)

	if len(s.Price) > 0 {
		for _, sPrice := range s.Price {
			if sPrice.PriceTypeCode > 0 {
				xmlElementSupplyDetailPrice(id, s.SupplierName, &sPrice)
			}
		}
	}
}

func xmlElementSupplyDetailPrice(id string, supplierName string, p *Price) {
	createTable(p)
	insertStmt := getInsertStmt(p)
	_, stmtErr := insertStmt.Exec(
		id,
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
	createTable(m)
	insertStmt := getInsertStmt(m)
	_, stmtErr := insertStmt.Exec(
		id,
		m.AgentName,
		m.AgentRole,
		m.MarketCountry,
		m.MarketPublishingStatus)
	handleErr(stmtErr)
}
