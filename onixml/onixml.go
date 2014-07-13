package onixml

// An example streaming XML parser.
// Initial Source: https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

import (
	"encoding/xml"
	"fmt"
	"database/sql"
	"net/url"
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

type ProductIdentifier struct {
	ProductIDType int   `xml:"ProductIDType",sql:"bigint(14)"`
	IDValue       string `xml:"IDValue"`
}

type Title struct {
	TitleType int    `xml:"Title>TitleType"`
	TitleText string `xml:"Title>TitleText"`
}

type Website struct {
	WebsiteLink string `xml:"Website>WebsiteLink"`
}

type Contributor struct {
	SequenceNumber     int    `xml:"Contributor>SequenceNumber"`
	ContributorRole    string `xml:"Contributor>ContributorRole"`
	PersonNameInverted string `xml:"Contributor>PersonNameInverted"`
	KeyNames           string `xml:"Contributor>KeyNames"`
}

type Extent struct {
	ExtentType      int    `xml:"Extent>ExtentType",sql"int(10) NOT NULL DEFAULT 0"`
	ExtentValue     int    `xml:"Extent>ExtentValue",sql"int(10) NOT NULL DEFAULT 0"`
	ExtentUnit      int    `xml:"Extent>ExtentUnit",sql"int(10) NOT NULL DEFAULT 0"`
}

type OtherText struct {
	TextTypeCode int    `xml:"OtherText>TextTypeCode",sql"int(10) NOT NULL"`
	Text         string   `xml:"OtherText>Text",sql"text NULL"`
}

type MediaFile struct {
	MediaFileTypeCode     int    `xml:"MediaFile>MediaFileTypeCode"`
	MediaFileLinkTypeCode int    `xml:"MediaFile>MediaFileLinkTypeCode"`
	MediaFileLink         string `xml:"MediaFile>MediaFileLink"`
}

type Imprint struct {
	ImprintName string `xml:"Imprint>ImprintName"`
}
type Publisher struct {
	PublishingRole int    `xml:"Publisher>PublishingRole"`
	PublisherName  string `xml:"Publisher>PublisherName"`
}

type SalesRights struct {
	SalesRightsType int    `xml:"SalesRights>SalesRightsType",sql:"int(10) NOT NULL"`
	RightsCountry  string `xml:"SalesRights>RightsCountry",sql:"varchar(2) NULL"`
}

type Measure struct {
	MeasureTypeCode int    `xml:"Measure>MeasureTypeCode"`
	Measurement     int    `xml:"Measure>Measurement"`
	MeasureUnitCode string `xml:"Measure>MeasureUnitCode"`
}

type RelatedProduct struct {
	RelationCode  int    `xml:"RelatedProduct>RelationCode",sql:"int(10) NOT NULL"`
	ProductIDType int    `xml:"RelatedProduct>ProductIdentifier>ProductIDType",sql:"int(10) NOT NULL"`
	IDValue       string    `xml:"RelatedProduct>ProductIdentifier>IDValue",sql:"bigint(15) NOT NULL"`
}

type SupplyDetail struct {
	SupplierName        string  `xml:"SupplyDetail>SupplierName"`
	SupplierRole        int     `xml:"SupplyDetail>SupplierRole"`
	SupplyToCountry     string  `xml:"SupplyDetail>SupplyToCountry"`
	ProductAvailability int     `xml:"SupplyDetail>ProductAvailability"`
	OnHand              int     `xml:"SupplyDetail>Stock>OnHand"`
	OnOrder             int     `xml:"SupplyDetail>Stock>OnOrder"`
	PackQuantity        int     `xml:"SupplyDetail>PackQuantity"`
	PriceTypeCode       int     `xml:"SupplyDetail>Price>PriceTypeCode"`
	PriceAmount         float32 `xml:"SupplyDetail>Price>PriceAmount"`
	CurrencyCode        string  `xml:"SupplyDetail>Price>CurrencyCode"`
	CountryCode         string  `xml:"SupplyDetail>Price>CountryCode"`
}

type MarketRepresentation struct {
	AgentName              string `xml:"MarketRepresentation>AgentName"`
	AgentRole              string `xml:"MarketRepresentation>AgentRole"`
	MarketCountry          string `xml:"MarketRepresentation>MarketCountry"`
	MarketPublishingStatus int    `xml:"MarketRepresentation>MarketPublishingStatus"`
}

type Product struct {
	RecordReference   string `xml:"RecordReference",sql:"bigint(15) NOT NULL DEFAULT 0"`
	NotificationType  int `xml:"NotificationType"`
	ProductIdentifier []ProductIdentifier
	ProductForm       string `xml:"ProductForm"`
	ProductFormDetail string `xml:"ProductFormDetail"`
	Title
	Website
	Contributor
	Extent
	EditionNumber  int    `xml:"EditionNumber"`
	NumberOfPages  int    `xml:"NumberOfPages"`
	BICMainSubject string `xml:"BICMainSubject"`
	OtherText
	AudienceCode   int    `xml:"AudienceCode"`
	MediaFile
	Imprint
	Publisher
	SalesRights
	PublishingStatus   int    `xml:"PublishingStatus"`
	PublicationDate    string `xml:"PublicationDate"`
	YearFirstPublished string `xml:"YearFirstPublished"`
	Measure
	RelatedProduct
	SupplyDetail
	MarketRepresentation
}

func CanonicalizeTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
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
				if "" != prod.Website.WebsiteLink {
					xmlElementWebsite(prod.RecordReference, &prod.Website)
				}
				if prod.Contributor.SequenceNumber > 0 {
					xmlElementContributor(prod.RecordReference, &prod.Contributor)
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
	_, stmtErr := insertStmt.Exec(
		// avoiding reflection
		prod.RecordReference,
		prod.RecordReference,
		prod.NotificationType,
		prod.ProductForm,
		prod.ProductFormDetail,
		prod.EditionNumber,
		prod.NumberOfPages,
		prod.BICMainSubject,
		prod.AudienceCode,
		prod.PublishingStatus,
		prod.PublicationDate,
		prod.YearFirstPublished)
	handleErr(stmtErr)

}

func xmlElementProductIdentifier(id string, prodIdentifier *ProductIdentifier) {
	createTable(prodIdentifier)
	insertStmt := getInsertStmt(prodIdentifier)

	_, stmtErr := insertStmt.Exec(
		id,
		prodIdentifier.ProductIDType,
		prodIdentifier.IDValue)
	handleErr(stmtErr)
}

func xmlElementTitle(id string, title *Title) {
	createTable(title)
	insertStmt := getInsertStmt(title)
	_, stmtErr := insertStmt.Exec(
		id,
		title.TitleType,
		title.TitleText)
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
