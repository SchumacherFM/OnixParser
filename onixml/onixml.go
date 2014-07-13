package onixml

// An example streaming XML parser.
// Initial Source: https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

import (
	"encoding/xml"
	"fmt"
	"database/sql"
	"os"
	"reflect"
	"strings"
	"../sqlCreator"
)

type tableColumn struct {
	Table, Column string
}

var (
	dbCon *sql.DB
	tablePrefix string
	preparedInsertStmt = make(map[string]*sql.Stmt)
)

func SetConnection(aCon *sql.DB) {
	dbCon = aCon
}
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

func handleErr(theErr error) {
	if nil != theErr {
		panic(theErr.Error())
	}
}

func OnixmlDecode(inputFile string) int {
	sqlCreator.SetTablePrefix(tablePrefix)
	total := 0
	xmlFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return total
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, dtErr := decoder.Token()
		if t == nil {
			break
		}
		handleErr(dtErr)

		// Inspect the type of the token just read.
		switch se := t.(type) { // wtf? I don't understand this magic
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "Product"
			if inElement == "Product" {
				var prod Product
				// decode a whole chunk of following XML into the
				// variable prod which is a Product (se above)
				decErr := decoder.DecodeElement(&prod, &se)
				handleErr(decErr)
				xmlElementProduct(&prod)
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
				if prod.Extent.ExtentType > 0 {
					xmlElementExtent(prod.RecordReference, &prod.Extent)
				}
				if prod.OtherText.TextTypeCode > 0 {
					xmlElementOtherText(prod.RecordReference, &prod.OtherText)
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
				if "" != prod.SupplyDetail.SupplierName {
					xmlElementSupplyDetail(prod.RecordReference, &prod.SupplyDetail)
				}
				if "" != prod.MarketRepresentation.AgentName {
					xmlElementMarketRepresentation(prod.RecordReference, &prod.MarketRepresentation)
				}


				total++
			}
		default:
		}
	}
	return total
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
		t.TitleText)
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
		c.KeyNames)
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
		s.OnHand,
		s.OnOrder,
		s.PackQuantity,
		s.PriceTypeCode,
		s.PriceAmount,
		s.CurrencyCode,
		s.CountryCode)
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
