package onixml

// An example streaming XML parser.
// Initial Source: https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

import (
	"encoding/xml"
	"fmt"
	//"log"
	"database/sql"
	"net/url"
	"os"
	"strings"
)

var (
	dbCon *sql.DB
)

func SetConnection(aCon *sql.DB) {
	dbCon = aCon
}

// Here is an example article from the Wikipedia XML dump
//
// <page>
// 	<title>Apollo 11</title>
//      <redirect title="Foo bar" />
// 	...
// 	<revision>
// 	...
// 	  <text xml:space="preserve">
// 	  {{Infobox Space mission
// 	  |mission_name=&lt;!--See above--&gt;
// 	  |insignia=Apollo_11_insignia.png
// 	...
// 	  </text>
// 	</revision>
// </page>
//
// Note how the tags on the fields of Page and Redirect below
// describe the XML schema structure.

func handleErr(theErr error) {
	if nil != theErr {
		panic(theErr.Error())
	}
}

type ProductIdentifier struct {
	ProductIDType int   `xml:"ProductIDType"`
	IDValue       int64 `xml:"IDValue"`
}
type Title struct {
	TitleType int    `xml:"TitleType"`
	TitleText string `xml:"TitleText"`
}

type Website struct {
	WebsiteLink string `xml:"WebsiteLink"`
}

type Contributor struct {
	SequenceNumber     int    `xml:"SequenceNumber"`
	ContributorRole    string `xml:"ContributorRole"`
	PersonNameInverted string `xml:"PersonNameInverted"`
	KeyNames           string `xml:"KeyNames"`
}

type MediaFile struct {
	MediaFileTypeCode     int    `xml:"MediaFileTypeCode"`
	MediaFileLinkTypeCode int    `xml:"MediaFileLinkTypeCode"`
	MediaFileLink         string `xml:"MediaFileLink"`
}

type Imprint struct {
	ImprintName string `xml:"ImprintName"`
}
type Publisher struct {
	PublishingRole int    `xml:"PublishingRole"`
	PublisherName  string `xml:"PublisherName"`
}

type Measure struct {
	MeasureTypeCode int    `xml:"MeasureTypeCode"`
	Measurement     int    `xml:"Measurement"`
	MeasureUnitCode string `xml:"MeasureUnitCode"`
}

type SupplyDetail struct {
	SupplierName        string  `xml:"SupplierName"`
	SupplierRole        int     `xml:"SupplierRole"`
	SupplyToCountry     string  `xml:"SupplyToCountry"`
	ProductAvailability int     `xml:"ProductAvailability"`
	OnHand              int     `xml:"Stock>OnHand"`
	OnOrder             int     `xml:"Stock>OnOrder"`
	PackQuantity        int     `xml:"PackQuantity"`
	PriceTypeCode       int     `xml:"Price>PriceTypeCode"`
	PriceAmount         float32 `xml:"Price>PriceAmount"`
	CurrencyCode        string  `xml:"Price>CurrencyCode"`
	CountryCode         string  `xml:"Price>CountryCode"`
}

type MarketRepresentation struct {
	AgentName              string `xml:"AgentName"`
	AgentRole              string `xml:"AgentRole"`
	MarketCountry          string `xml:"MarketCountry"`
	MarketPublishingStatus int    `xml:"MarketPublishingStatus"`
}

type Product struct {
	RecordReference   int64 `xml:"RecordReference"`
	NotificationType  int `xml:"NotificationType"`
	ProductIdentifier []ProductIdentifier
	ProductForm       string `xml:"ProductForm"`
	ProductFormDetail string `xml:"ProductFormDetail"`
	Title
	Website
	Contributor
	EditionNumber  int    `xml:"EditionNumber"`
	NumberOfPages  int    `xml:"NumberOfPages"`
	BICMainSubject string `xml:"BICMainSubject"`
	AudienceCode   int    `xml:"AudienceCode"`
	MediaFile
	Imprint
	Publisher
	PublishingStatus   int    `xml:"PublishingStatus"`
	PublicationDate    string `xml:"PublicationDate"`
	YearFirstPublished string `xml:"YearFirstPublished"`
	Measure
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
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) { // wtf? I don't understand this magic
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "Product"
			if inElement == "Product" {
				var prod Product
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)

				decoder.DecodeElement(&prod, &se)
				//fmt.Printf("%v\n", prod)
				for pidI, pidE := range prod.ProductIdentifier {
					fmt.Printf("%d: %d => %d\n", prod.RecordReference, pidI, pidE)
				}


				// Do some stuff with the page.
				//				p.Title = CanonicalizeTitle(p.Title)
				//					m := filter.MatchString(p.Title)
				//					if !m && p.Redir.Title == "" {
				//						WritePage(p.Title, p.Text)
				//						total++
				//					}
			}
		default:
		}

	}
	return total
}
