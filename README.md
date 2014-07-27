# OnixParser

A parser for the [EDItEUR ONIX xml file format](http://www.editeur.org/83/Overview/) written in GoLang.

My first learning project. Of course it contains many *mistakes* but the binary runs quite nice and at Zookal.com we're
using it in production. If you find an error or incorrect usage of Go, let me know :-)

Inspired for learning by [http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)

# XML Data

[http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd](http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd)

### Test data Onix Data Feed

You can download test data from [http://www.oup.com.au/help_and_advice/booksellers](http://www.oup.com.au/help_and_advice/booksellers)

The Complete File or The Incremental File.

Not all XML elements are matched because structure in DTD is unclear and we don't need all elements. If one element is
missing send me an email or a pull request.

# Performance

Parsing a **3GB XML** file with **998673 products** needs **14m36.525541544s** to import that data **into MySQL 5.5**.

Hardware used:

- MacBook Air
- 13-inch, Mid 2012
- Processor  1.8 GHz Intel Core i5
- Memory  8 GB 1600 MHz DDR3
- APPLE SATA SSD SM256E Media

There are several options on the command line:

```
$ go run OnixParser.go -h
OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm
This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html
Usage of OnixParser:
  -db="test": MySQL db name
  -host="127.0.0.1": MySQL host name
  -infile="": Input file path
  -logfile="": Logfile name, if empty direct output
  -moc=20: Max MySQL open connections
  -outdir="": Dir for CSV output file for reading into MySQL, if empty writes to /tmp/
  -pass="test": MySQL password
  -tablePrefix="gonix_": Table name prefix
  -user="test": MySQL user name
  -v=false: Increase verbosity
exit status 2
```

The speed of import can only be achieved by using `LOAD DATA INFILE LOCAL` with MySQL. So check your config if the option
has been enabled.

# License

General Public License

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

# Author

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

Made in Sydney, Australia :-)

If you consider a donation please contribute to: [http://www.seashepherd.org/](http://www.seashepherd.org/)
