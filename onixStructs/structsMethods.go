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
package onixStructs

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

var (
	internalByteCounter     = make(map[string]int)
	currentWriteToTableName = make(map[string]string)
)

func getCurrentWriteToTableName(tableName string) string {
	tn, isSet := currentWriteToTableName[tableName]
	if false == isSet {
		return tableName
	}
	return tn
}

// get around of mysql max allowed packet which is hardcoded in the mysql driver at 8MB :-(
func countByte(tableName string, bytes int) {
	counted, isSet := internalByteCounter[tableName]
	if false == isSet {
		internalByteCounter[tableName] = 0
	}
	internalByteCounter[tableName] = bytes + counted
}

func moreThanMySqlMaxAllowedPacket(tableName string) bool {
	if internalByteCounter[tableName] > appConfig.MaxPacketSize {
		internalByteCounter[tableName] = 0
		return true
	}
	return false
}

func writeOneElementToFile(anyStruct interface{}, args map[int]string) (int, error) {
	tableName := appConfig.GetNameOfStruct(anyStruct)
	mapLen := len(args) - 1
	var buffer bytes.Buffer

	// important to keep the correct order of the map
	keys := []int{}
	for k := range args {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		buffer.WriteByte(appConfig.Csv.Enclosure)
		buffer.WriteString(args[k])
		buffer.WriteByte(appConfig.Csv.Enclosure)
		if k < mapLen {
			buffer.WriteByte(appConfig.Csv.Delimiter)
		}
	}
	buffer.WriteByte(appConfig.Csv.LineEnding)
	countByte(tableName, buffer.Len())

	if true == moreThanMySqlMaxAllowedPacket(tableName) {
		// create new file
		nextTableName := appConfig.GetNextTableName(tableName)
		currentWriteToTableName[tableName] = nextTableName
	}
	writeToTn := getCurrentWriteToTableName(tableName)
	return appConfig.WriteBytes(writeToTn, buffer.Bytes())
}

func (p *Product) Xml2CsvRoot() {
	_, writeErr := writeOneElementToFile(p, map[int]string{
		0:  p.RecordReference,
		1:  p.RecordReference,
		2:  strconv.Itoa(p.NotificationType),
		3:  p.ProductForm,
		4:  p.ProductFormDetail,
		5:  p.EditionNumber,
		6:  strings.Replace(p.NumberOfPages, ",", "", -1),
		7:  p.IllustrationsNote,
		8:  p.BICMainSubject,
		9:  strings.TrimLeft(strings.TrimSpace(p.AudienceCode), "0"),
		10: strconv.Itoa(p.PublishingStatus),
		11: p.PublicationDate,
		12: p.YearFirstPublished,
	})
	appConfig.HandleErr(writeErr)
}

func (p *ProductIdentifier) Xml2Csv(id string) {
	if p.ProductIDType > 0 {
		_, writeErr := writeOneElementToFile(p, map[int]string{
			0: id,
			1: strconv.Itoa(p.ProductIDType),
			2: p.IDValue,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (t *Title) Xml2Csv(id string) {
	if t.TitleType > 0 {
		_, writeErr := writeOneElementToFile(t, map[int]string{
			0: id,
			1: strconv.Itoa(t.TitleType),
			2: t.TitleText,
			3: t.TitlePrefix,
			4: t.TitleWithoutPrefix,
		})
		appConfig.HandleErr(writeErr)
	}
}
func (s *Series) Xml2Csv(id string) {
	if "" != s.TitleOfSeries || "" != s.NumberWithinSeries {
		_, writeErr := writeOneElementToFile(s, map[int]string{
			0: id,
			1: s.TitleOfSeries,
			2: s.NumberWithinSeries,
		})
		appConfig.HandleErr(writeErr)
	}
}
func (w *Website) Xml2Csv(id string) {
	if "" != w.WebsiteLink {
		_, writeErr := writeOneElementToFile(w, map[int]string{
			0: id,
			1: w.WebsiteLink,
		})
		appConfig.HandleErr(writeErr)
	}
}
func (c *Contributor) Xml2Csv(id string) {
	if c.SequenceNumber > 0 {
		_, writeErr := writeOneElementToFile(c, map[int]string{
			0: id,
			1: strconv.Itoa(c.SequenceNumber),
			2: c.ContributorRole,
			3: c.PersonNameInverted,
			4: c.TitlesBeforeNames,
			5: c.KeyNames,
		})
		appConfig.HandleErr(writeErr)
	}
}
func (s *Subject) Xml2Csv(id string) {
	if s.SubjectSchemeIdentifier > 0 {
		_, writeErr := writeOneElementToFile(s, map[int]string{
			0: id,
			1: strconv.Itoa(s.SubjectSchemeIdentifier),
			2: s.SubjectCode,
		})
		appConfig.HandleErr(writeErr)
	}
}
func (e *Extent) Xml2Csv(id string) {
	if e.ExtentType > 0 {
		_, writeErr := writeOneElementToFile(e, map[int]string{
			0: id,
			1: strconv.Itoa(e.ExtentType),
			2: strconv.Itoa(e.ExtentValue),
			3: strconv.Itoa(e.ExtentUnit),
		})
		appConfig.HandleErr(writeErr)
	}
}
func (s *SupplyDetail) Xml2Csv(id string) {
	if "" != s.SupplierName {
		_, writeErr := writeOneElementToFile(s, map[int]string{
			0: id,
			1: s.SupplierName,
			2: strconv.Itoa(s.SupplierRole),
			3: s.SupplyToCountry,
			4: strconv.Itoa(s.ProductAvailability),
			5: s.ExpectedShipDate,
			6: strconv.Itoa(s.OnHand),
			7: strconv.Itoa(s.OnOrder),
			8: strconv.Itoa(s.PackQuantity),
		})
		appConfig.HandleErr(writeErr)

		if len(s.Price) > 0 {
			for _, sPrice := range s.Price {
				if sPrice.PriceTypeCode > 0 {
					sPrice.Xml2Csv(id, s.SupplierName)
				}
			}
		}
	}
}

func (p *Price) Xml2Csv(id string, supplierName string) {
	_, writeErr := writeOneElementToFile(p, map[int]string{
		0: id,
		1: supplierName,
		2: strconv.Itoa(p.PriceTypeCode),
		3: strconv.Itoa(p.DiscountCodeType),
		4: p.DiscountCode,
		5: p.PriceAmount,
		6: p.CurrencyCode,
		7: p.CountryCode,
	})
	appConfig.HandleErr(writeErr)
}

func (o *OtherText) Xml2Csv(id string) {
	if o.TextTypeCode > 0 {
		_, writeErr := writeOneElementToFile(o, map[int]string{
			0: id,
			1: strconv.Itoa(o.TextTypeCode),
			2: o.Text,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (m *MediaFile) Xml2Csv(id string) {
	if m.MediaFileTypeCode > 0 {
		_, writeErr := writeOneElementToFile(m, map[int]string{
			0: id,
			1: strconv.Itoa(m.MediaFileTypeCode),
			2: strconv.Itoa(m.MediaFileLinkTypeCode),
			3: m.MediaFileLink,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (i *Imprint) Xml2Csv(id string) {
	if "" != i.ImprintName {
		_, writeErr := writeOneElementToFile(i, map[int]string{
			0: id,
			1: i.ImprintName,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (p *Publisher) Xml2Csv(id string) {
	if p.PublishingRole > 0 {
		_, writeErr := writeOneElementToFile(p, map[int]string{
			0: id,
			1: strconv.Itoa(p.PublishingRole),
			2: p.PublisherName,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (s *SalesRights) Xml2Csv(id string) {
	if s.SalesRightsType > 0 {
		_, writeErr := writeOneElementToFile(s, map[int]string{
			0: id,
			1: strconv.Itoa(s.SalesRightsType),
			2: s.RightsCountry,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (s *SalesRestriction) Xml2Csv(id string) {
	if s.SalesRestrictionType > 0 {
		_, writeErr := writeOneElementToFile(s, map[int]string{
			0: id,
			1: strconv.Itoa(s.SalesRestrictionType),
		})
		appConfig.HandleErr(writeErr)
	}
}

func (m *Measure) Xml2Csv(id string) {
	if m.MeasureTypeCode > 0 {
		_, writeErr := writeOneElementToFile(m, map[int]string{
			0: id,
			1: strconv.Itoa(m.MeasureTypeCode),
			2: m.Measurement,
			3: m.MeasureUnitCode,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (r *RelatedProduct) Xml2Csv(id string) {
	if r.ProductIDType > 0 {
		_, writeErr := writeOneElementToFile(r, map[int]string{
			0: id,
			1: strconv.Itoa(r.RelationCode),
			2: strconv.Itoa(r.ProductIDType),
			3: r.IDValue,
		})
		appConfig.HandleErr(writeErr)
	}
}

func (m *MarketRepresentation) Xml2Csv(id string) {
	if "" != m.AgentName {
		_, writeErr := writeOneElementToFile(m, map[int]string{
			0: id,
			1: m.AgentName,
			2: strconv.Itoa(m.AgentRole),
			3: m.MarketCountry,
			4: strconv.Itoa(m.MarketPublishingStatus),
		})
		appConfig.HandleErr(writeErr)
	}
}
