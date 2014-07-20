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
	"strings"
)

func (p *Product) writeToDb(id string) {
	iSql := getInsertStmt(p)
	// _, stmtErr := insertStmt.Exec.Call(p) => that would be nice ... but how?
	// static typed language and that would cost performance
	/* sometimes number can be 1,234 */
	// avoiding reflection
	_, stmtErr := appConfig.dbCon.Exec(
		iSql,
		p.RecordReference,
		p.RecordReference,
		p.NotificationType,
		p.ProductForm,
		p.ProductFormDetail,
		p.EditionNumber,
		strings.Replace(p.NumberOfPages, ",", "", -1),
		p.IllustrationsNote,
		p.BICMainSubject,
		p.AudienceCode,
		p.PublishingStatus,
		p.PublicationDate,
		p.YearFirstPublished)
	handleErr(stmtErr)
}

func (p *ProductIdentifier) writeToDb(id string) {
	iSql := getInsertStmt(p)
	_, stmtErr := appConfig.dbCon.Exec(
		iSql,
		id,
		p.ProductIDType,
		p.IDValue)
	handleErr(stmtErr)
}

func (t *Title) writeToDb(id string) {
	if t.TitleType > 0 {
		iSql := getInsertStmt(t)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql, id,
			t.TitleType,
			t.TitleText,
			t.TitlePrefix,
			t.TitleWithoutPrefix)
		handleErr(stmtErr)
	}
}
func (s *Series) writeToDb(id string) {
	if "" != s.TitleOfSeries || "" != s.NumberWithinSeries {
		iSql := getInsertStmt(s)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql,
			id,
			s.TitleOfSeries,
			s.NumberWithinSeries)
		handleErr(stmtErr)
	}
}
func (w *Website) writeToDb(id string) {
	if "" != w.WebsiteLink {
		iSql := getInsertStmt(w)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql,
			id,
			w.WebsiteLink)
		handleErr(stmtErr)
	}
}
func (c *Contributor) writeToDb(id string) {
	if c.SequenceNumber > 0 {
		iSql := getInsertStmt(c)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql,
			id,
			c.SequenceNumber,
			c.ContributorRole,
			c.PersonNameInverted,
			c.TitlesBeforeNames,
			c.KeyNames)
		handleErr(stmtErr)
	}
}
func (s *Subject) writeToDb(id string) {
	if s.SubjectSchemeIdentifier > 0 {
		iSql := getInsertStmt(s)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql,
			id,
			s.SubjectSchemeIdentifier,
			s.SubjectCode)
		handleErr(stmtErr)
	}
}
func (e *Extent) writeToDb(id string) {
	if e.ExtentType > 0 {
		iSql := getInsertStmt(e)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql, id,
			e.ExtentType,
			e.ExtentValue,
			e.ExtentUnit)
		handleErr(stmtErr)
	}
}
func (s *SupplyDetail) writeToDb(id string) {
	if "" != s.SupplierName {
		iSql := getInsertStmt(s)
		_, stmtErr := appConfig.dbCon.Exec(
			iSql, id,
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
}
