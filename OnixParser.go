package main

// An example streaming XML parser.
// Initial Source: https://github.com/dps/go-xml-parse/blob/master/go-xml-parse.go

import (
	"bufio"
	//	"encoding/xml"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	timeStart   = time.Now()
	inputFile   = flag.String("infile", "demo-availability.xml", "Input file path")
	dbHost      = flag.String("host", "127.0.0.1", "MySQL host name")
	dbDb        = flag.String("db", "test", "MySQL db name")
	dbUser      = flag.String("user", "test", "MySQL user name")
	dbPass      = flag.String("pass", "test", "MySQL password")
	tablePrefix = flag.String("tablePrefix", "gonix_", "Table name prefix")
	tableColumns = make([][]string,50)
	dbCon *sql.DB
)
var filter, _ = regexp.Compile("^file:.*|^talk:.*|^special:.*|^wikipedia:.*|^wiktionary:.*|^user:.*|^user_talk:.*")

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

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title string   `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text  string   `xml:"revision>text"`
}

func CanonicalizeTitle(title string) string {
	can := strings.ToLower(title)
	can = strings.Replace(can, " ", "_", -1)
	can = url.QueryEscape(can)
	return can
}

func WritePage(title string, text string) {
	outFile, err := os.Create("out/docs/" + title)
	if err == nil {
		writer := bufio.NewWriter(outFile)
		defer outFile.Close()
		writer.WriteString(text)
		writer.Flush()
	}
}

func initDatabase() {
	dbCon, dbConErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", *dbUser, *dbPass, *dbHost, *dbDb))
	if dbConErr != nil {
		panic(dbConErr.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer dbCon.Close()

	// Open doesn't open a connection. Validate DSN data:
	err := dbCon.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// delete already created tables
	// @todo go on here
	rows, err := dbCon.Query("SHOW TABLES FROM ?", dbDb)
	if err != nil { /* error handling */}
	partages := make([]*Partage, 0, 10)
	var ida, idb uint
	for rows.Next() {
		err = rows.Scan(&ida, &idb)
		if err != nil { /* error handling */}
		partages = append(partages, &Partage{ida, idb})
	}
}

func printDuration() {
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	fmt.Printf("XML Parser took %dh %dm %fs to run.\n", int(duration.Hours()), int(duration.Minutes()), duration.Seconds())
	fmt.Printf("XML Parser took %v to run.\n", duration)
}

func main() {
	flag.Parse()
	initDatabase()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	//	decoder := xml.NewDecoder(xmlFile)
	total := 0
	//	var inElement string
	//	for {
	//		// Read tokens from the XML document in a stream.
	//		t, _ := decoder.Token()
	//		if t == nil {
	//			break
	//		}
	//		// Inspect the type of the token just read.
	//		switch se := t.(type) {
	//		case xml.StartElement:
	//			// If we just read a StartElement token
	//			inElement = se.Name.Local
	//			// ...and its name is "page"
	//			if inElement == "page" {
	//				var p Page
	//				// decode a whole chunk of following XML into the
	//				// variable p which is a Page (se above)
	//				decoder.DecodeElement(&p, &se)
	//
	//				// Do some stuff with the page.
	//				p.Title = CanonicalizeTitle(p.Title)
	//				m := filter.MatchString(p.Title)
	//				if !m && p.Redir.Title == "" {
	//					WritePage(p.Title, p.Text)
	//					total++
	//				}
	//			}
	//		default:
	//		}
	//
	//	}

	fmt.Printf("Total articles: %d \n", total)
	printDuration()
}
