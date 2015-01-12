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
	"sync"

	"github.com/SchumacherFM/OnixParser/gonfig"
)

var (
	appConfig *gonfig.AppConfiguration
	theSyncer = new(sync.Mutex)
)

func ParseXmlElementsConcurrent(prod *Product, appConfigArg *gonfig.AppConfiguration, wg *sync.WaitGroup) {
	// as we are in another thread set the dbCon new
	appConfig = appConfigArg
	defer wg.Done()
	theSyncer.Lock()
	defer theSyncer.Unlock() // sync between the concurrent processes for writing to the files

	prod.Xml2CsvRoot()

	if len(prod.ProductIdentifier) > 0 {
		for _, prodIdentifier := range prod.ProductIdentifier {
			prodIdentifier.Xml2Csv(prod.RecordReference)
		}
	}
	prod.Title.Xml2Csv(prod.RecordReference)
	prod.Series.Xml2Csv(prod.RecordReference)
	prod.Website.Xml2Csv(prod.RecordReference)
	prod.Extent.Xml2Csv(prod.RecordReference)

	if len(prod.Contributor) > 0 {
		for _, prodContributor := range prod.Contributor {
			prodContributor.Xml2Csv(prod.RecordReference)
		}
	}

	if len(prod.Subject) > 0 {
		for _, prodSubject := range prod.Subject {
			prodSubject.Xml2Csv(prod.RecordReference)
		}
	}
	if len(prod.SupplyDetail) > 0 {
		for _, prodSupplyDetail := range prod.SupplyDetail {
			prodSupplyDetail.Xml2Csv(prod.RecordReference)
		}
	}

	if len(prod.OtherText) > 0 {
		for _, prodOtherText := range prod.OtherText {
			prodOtherText.Xml2Csv(prod.RecordReference)
		}
	}
	prod.MediaFile.Xml2Csv(prod.RecordReference)
	prod.Imprint.Xml2Csv(prod.RecordReference)
	prod.Publisher.Xml2Csv(prod.RecordReference)
	prod.SalesRights.Xml2Csv(prod.RecordReference)
	prod.SalesRestriction.Xml2Csv(prod.RecordReference)
	prod.Measure.Xml2Csv(prod.RecordReference)
	prod.RelatedProduct.Xml2Csv(prod.RecordReference)
	prod.MarketRepresentation.Xml2Csv(prod.RecordReference)

}
