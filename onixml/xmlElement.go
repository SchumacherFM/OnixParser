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
	"sync" // for concurrency
	"github.com/SchumacherFM/OnixParser/gonfig"
)

func parseXmlElementsConcurrent(prod *Product, appConfig *gonfig.AppConfiguration, wg *sync.WaitGroup) {
	// as we are in another thread set the dbCon new
	SetAppConfig(appConfig)
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
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		o.TextTypeCode,
		o.Text)
	appConfig.HandleErr(stmtErr)
}

func xmlElementMediaFile(id string, m *MediaFile) {
	iSql := getInsertStmt(m)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		m.MediaFileTypeCode,
		m.MediaFileLinkTypeCode,
		m.MediaFileLink)
	appConfig.HandleErr(stmtErr)
}

func xmlElementImprint(id string, i *Imprint) {
	iSql := getInsertStmt(i)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		i.ImprintName)
	appConfig.HandleErr(stmtErr)
}

func xmlElementPublisher(id string, p *Publisher) {
	iSql := getInsertStmt(p)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		p.PublishingRole,
		p.PublisherName)
	appConfig.HandleErr(stmtErr)
}
func xmlElementSalesRights(id string, s *SalesRights) {
	iSql := getInsertStmt(s)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		s.SalesRightsType,
		s.RightsCountry)
	appConfig.HandleErr(stmtErr)
}
func xmlElementSalesRestriction(id string, s *SalesRestriction) {
	iSql := getInsertStmt(s)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		s.SalesRestrictionType)
	appConfig.HandleErr(stmtErr)
}

func xmlElementMeasure(id string, m *Measure) {
	iSql := getInsertStmt(m)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		m.MeasureTypeCode,
		m.Measurement,
		m.MeasureUnitCode)
	appConfig.HandleErr(stmtErr)
}

func xmlElementRelatedProduct(id string, r *RelatedProduct) {
	iSql := getInsertStmt(r)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		r.RelationCode,
		r.ProductIDType,
		r.IDValue)
	appConfig.HandleErr(stmtErr)
}

func xmlElementSupplyDetailPrice(id string, supplierName string, p *Price) {
	iSql := getInsertStmt(p)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		supplierName,
		p.PriceTypeCode,
		p.DiscountCodeType,
		p.DiscountCode,
		p.PriceAmount,
		p.CurrencyCode,
		p.CountryCode)
	appConfig.HandleErr(stmtErr)
}

func xmlElementMarketRepresentation(id string, m *MarketRepresentation) {
	iSql := getInsertStmt(m)
	_, stmtErr := appConfig.GetConnection().Exec(
		iSql, id,
		m.AgentName,
		m.AgentRole,
		m.MarketCountry,
		m.MarketPublishingStatus)
	appConfig.HandleErr(stmtErr)
}
