package onixml

// An example streaming XML parser.
// Initial Source: https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go
type (
	ProductIdentifier struct {
		ProductIDType int    `xml:"ProductIDType" sql:"bigint(14)"`
		IDValue       string `xml:"IDValue" sql:"varchar(255) NULL"`
	}

	Title struct {
		TitleType int    `xml:"Title>TitleType" sql:"int(10) NOT NULL"`
		TitleText string `xml:"Title>TitleText" sql:"varchar(255) NULL"`
	}
	Series struct {
		TitleOfSeries      string `xml:"Series>TitleOfSeries" sql:"varchar(255) NULL"`
		NumberWithinSeries string `xml:"Series>NumberWithinSeries" sql:"int(10) NOT NULL DEFAULT 0"`
	}

	Website struct {
		WebsiteLink string `xml:"Website>WebsiteLink" sql:"text NULL"`
	}

	Contributor struct {
		SequenceNumber     int    `xml:"SequenceNumber" sql:"int(10) NOT NULL"`
		ContributorRole    string `xml:"ContributorRole" sql:"varchar(255) NULL"`
		PersonNameInverted string `xml:"PersonNameInverted" sql:"varchar(255) NULL"`
		KeyNames           string `xml:"KeyNames" sql:"varchar(255) NULL"`
	}

	Subject struct {
		SubjectSchemeIdentifier int    `xml:"SubjectSchemeIdentifier" sql:"int(10) NOT NULL"`
		SubjectCode             string `xml:"SubjectCode" sql:"varchar(255) NULL"`
	}

	Extent struct {
		ExtentType  int `xml:"Extent>ExtentType" sql:"int(10) NOT NULL DEFAULT 0"`
		ExtentValue int `xml:"Extent>ExtentValue" sql:"int(10) NOT NULL DEFAULT 0"`
		ExtentUnit  int `xml:"Extent>ExtentUnit" sql:"int(10) NOT NULL DEFAULT 0"`
	}

	OtherText struct {
		TextTypeCode int    `xml:"TextTypeCode" sql:"int(10) NOT NULL"`
		Text         string `xml:"Text" sql:"text NULL"`
	}

	MediaFile struct {
		MediaFileTypeCode     int    `xml:"MediaFile>MediaFileTypeCode" sql:"int(10) NOT NULL"`
		MediaFileLinkTypeCode int    `xml:"MediaFile>MediaFileLinkTypeCode" sql:"int(10) NOT NULL"`
		MediaFileLink         string `xml:"MediaFile>MediaFileLink" sql:"text NULL"`
	}

	Imprint struct {
		ImprintName string `xml:"Imprint>ImprintName" sql:"varchar(255) NULL"`
	}
	Publisher struct {
		PublishingRole int    `xml:"Publisher>PublishingRole" sql:"int(10) NOT NULL"`
		PublisherName  string `xml:"Publisher>PublisherName" sql:"varchar(255) NULL"`
	}

	SalesRights struct {
		SalesRightsType int    `xml:"SalesRights>SalesRightsType" sql:"int(10) NOT NULL"`
		RightsCountry   string `xml:"SalesRights>RightsCountry" sql:"varchar(2) NULL"`
	}

	SalesRestriction struct {
		SalesRestrictionType int `xml:"SalesRestriction>SalesRestrictionType" sql:"int(10) NOT NULL"`
	}

	Measure struct {
		MeasureTypeCode int     `xml:"Measure>MeasureTypeCode" sql:"int(10) NOT NULL"`
		Measurement     float32 `xml:"Measure>Measurement" sql:"decimal(10,2) NOT NULL DEFAULT 0"`
		MeasureUnitCode string  `xml:"Measure>MeasureUnitCode" sql:"varchar(10) NULL"`
	}

	RelatedProduct struct {
		RelationCode  int    `xml:"RelatedProduct>RelationCode" sql:"int(10) NOT NULL"`
		ProductIDType int    `xml:"RelatedProduct>ProductIdentifier>ProductIDType" sql:"int(10) NOT NULL"`
		IDValue       string `xml:"RelatedProduct>ProductIdentifier>IDValue" sql:"bigint(15) NOT NULL"`
	}

	Price struct {
		SupplierName         string   `sql:"varchar(255) NOT NULL"` // only used in SQL table
		PriceTypeCode        int     `xml:"PriceTypeCode" sql:"int(10) NOT NULL DEFAULT 0"`
		DiscountCodeType     int     `xml:"DiscountCoded>DiscountCodeType" sql:"int(10) NOT NULL DEFAULT 0"`
		DiscountCode         string  `xml:"DiscountCoded>DiscountCode" sql:"varchar(10)  NULL"`
		PriceAmount          float32 `xml:"PriceAmount" sql:"decimal(10,2) NOT NULL DEFAULT 0"`
		CurrencyCode         string  `xml:"CurrencyCode" sql:"varchar(10) NULL"`
		CountryCode          string  `xml:"CountryCode" sql:"varchar(10) NULL"`
	}

	SupplyDetail struct {
		SupplierName         string  `xml:"SupplierName" sql:"varchar(255) NOT NULL"`
		SupplierRole         int     `xml:"SupplierRole" sql:"int(10) NOT NULL DEFAULT 0"`
		SupplyToCountry      string  `xml:"SupplyToCountry" sql:"varchar(255) NULL"`
		ProductAvailability  int     `xml:"ProductAvailability" sql:"int(10) NOT NULL DEFAULT 0"`
		ExpectedShipDate     string  `xml:"ExpectedShipDate" sql:"date NULL"`
		OnHand               int     `xml:"Stock>OnHand" sql:"int(10) NOT NULL DEFAULT 0"`
		OnOrder              int     `xml:"Stock>OnOrder" sql:"int(10) NOT NULL DEFAULT 0"`
		PackQuantity         int     `xml:"PackQuantity" sql:"int(10) NOT NULL DEFAULT 0"`
		Price                []Price
	}

	MarketRepresentation struct {
		AgentName              string `xml:"MarketRepresentation>AgentName" sql:"varchar(255) NOT NULL"`
		AgentRole              int    `xml:"MarketRepresentation>AgentRole" sql:"int(10) NOT NULL DEFAULT 0"`
		MarketCountry          string `xml:"MarketRepresentation>MarketCountry" sql:"varchar(4) NULL"`
		MarketPublishingStatus int    `xml:"MarketRepresentation>MarketPublishingStatus" sql:"int(10) NOT NULL DEFAULT 0"`
	}

	Product struct {
		RecordReference   string `xml:"RecordReference" sql:"bigint(15) NOT NULL DEFAULT 0"`
		NotificationType  int    `xml:"NotificationType" sql:"int(10) NOT NULL DEFAULT 0"`
		ProductIdentifier []ProductIdentifier
		ProductForm       string `xml:"ProductForm" sql:"varchar(20) NULL"`
		ProductFormDetail string `xml:"ProductFormDetail" sql:"varchar(20) NULL"`
		Series
		Title
		Website
		Contributor []Contributor
		Subject     []Subject
		Extent
		EditionNumber  string `xml:"EditionNumber" sql:"varchar(255) NULL"`
		NumberOfPages  string `xml:"NumberOfPages" sql:"int(10) NOT NULL DEFAULT 0"`
		BICMainSubject string `xml:"BICMainSubject" sql:"varchar(20) NULL"`
		OtherText      []OtherText
		AudienceCode   int `xml:"AudienceCode" sql:"int(10) NOT NULL DEFAULT 0"`
		MediaFile
		Imprint
		Publisher
		SalesRights
		SalesRestriction
		PublishingStatus   int    `xml:"PublishingStatus" sql:"int(10) NOT NULL DEFAULT 0"`
		PublicationDate    string `xml:"PublicationDate" sql:"varchar(255) NULL"`
		YearFirstPublished string `xml:"YearFirstPublished" sql:"varchar(255) NULL"`
		Measure
		RelatedProduct
		SupplyDetail    []SupplyDetail
		MarketRepresentation
	}
)
