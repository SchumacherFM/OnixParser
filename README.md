# OnixParser

This repo contains GoLang and PHP files to parse Onix xml files.

There are two branches:

- Master branch: slow synchronous read with ~5MB memory usage
- [Concurrency branch](https://github.com/SchumacherFM/OnixParser/tree/concurrency): Concurrent writes to the database. Memory usage depends on the file size. e.g. 278MB file requires 500MB RAM.

The concurrent execution is like a DDos attack to your database server. Be careful!

For more statistics please see the README in each branch.

The build binaries in each branch are for OSX 10.9 Darwin.

# XML Data

[http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd](http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd)

#### Test data Onix Data Feed

You can download test data from [http://www.oup.com.au/help_and_advice/booksellers](http://www.oup.com.au/help_and_advice/booksellers)

The Complete File or The Incremental File

# Go

Inspired for learning by [http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)

Not all XML elements are matched because structure in DTD is unclear and we don't need all elements.

XML Filesize: 278MB Output:

```
$ ./OnixParser --infile xmlFiles/oup_onix.xml --db test2 -v
OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm
This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html
2014/07/20 12:47:48 Dropped 15 existing tables
2014/07/20 12:47:56 7.983556469s for 1000 entities. Processed 1000, Mem alloc: 4.94MB
2014/07/20 12:48:05 8.141635991s for 1000 entities. Processed 2000, Mem alloc: 4.94MB
2014/07/20 12:48:13 8.101463417s for 1000 entities. Processed 3000, Mem alloc: 4.94MB
2014/07/20 12:48:20 7.773276568s for 1000 entities. Processed 4000, Mem alloc: 4.94MB
2014/07/20 12:48:28 7.958489608s for 1000 entities. Processed 5000, Mem alloc: 4.94MB
2014/07/20 12:48:37 8.314620045s for 1000 entities. Processed 6000, Mem alloc: 4.94MB
2014/07/20 12:48:45 8.113097432s for 1000 entities. Processed 7000, Mem alloc: 4.94MB
2014/07/20 12:48:53 8.146950940s for 1000 entities. Processed 8000, Mem alloc: 4.94MB
2014/07/20 12:49:00 7.492156122s for 1000 entities. Processed 9000, Mem alloc: 4.94MB
2014/07/20 12:49:07 6.941528537s for 1000 entities. Processed 10000, Mem alloc: 4.94MB
2014/07/20 12:49:15 7.202619983s for 1000 entities. Processed 11000, Mem alloc: 4.94MB
2014/07/20 12:49:22 7.331659125s for 1000 entities. Processed 12000, Mem alloc: 4.94MB
2014/07/20 12:49:30 7.757557180s for 1000 entities. Processed 13000, Mem alloc: 4.94MB
2014/07/20 12:49:37 7.247647428s for 1000 entities. Processed 14000, Mem alloc: 4.94MB
2014/07/20 12:49:45 8.044382828s for 1000 entities. Processed 15000, Mem alloc: 4.94MB
2014/07/20 12:49:53 7.581285830s for 1000 entities. Processed 16000, Mem alloc: 4.94MB
2014/07/20 12:50:00 7.566212804s for 1000 entities. Processed 17000, Mem alloc: 4.94MB
2014/07/20 12:50:07 7.348320872s for 1000 entities. Processed 18000, Mem alloc: 4.94MB
2014/07/20 12:50:15 7.773427457s for 1000 entities. Processed 19000, Mem alloc: 4.94MB
2014/07/20 12:50:23 7.453330394s for 1000 entities. Processed 20000, Mem alloc: 4.94MB
2014/07/20 12:50:30 7.695007291s for 1000 entities. Processed 21000, Mem alloc: 4.94MB
2014/07/20 12:50:38 7.684701404s for 1000 entities. Processed 22000, Mem alloc: 4.94MB
2014/07/20 12:50:46 7.629913776s for 1000 entities. Processed 23000, Mem alloc: 4.94MB
2014/07/20 12:50:53 7.290280485s for 1000 entities. Processed 24000, Mem alloc: 4.94MB
2014/07/20 12:50:59 6.475521434s for 1000 entities. Processed 25000, Mem alloc: 4.94MB
2014/07/20 12:51:07 7.600209449s for 1000 entities. Processed 26000, Mem alloc: 4.94MB
2014/07/20 12:51:15 7.453529103s for 1000 entities. Processed 27000, Mem alloc: 4.94MB
2014/07/20 12:51:22 7.025193358s for 1000 entities. Processed 28000, Mem alloc: 4.94MB
2014/07/20 12:51:29 7.507347268s for 1000 entities. Processed 29000, Mem alloc: 4.94MB
2014/07/20 12:51:36 7.302226818s for 1000 entities. Processed 30000, Mem alloc: 4.94MB
2014/07/20 12:51:44 7.577789023s for 1000 entities. Processed 31000, Mem alloc: 4.94MB
2014/07/20 12:51:51 6.966098416s for 1000 entities. Processed 32000, Mem alloc: 4.94MB
2014/07/20 12:51:59 7.733478405s for 1000 entities. Processed 33000, Mem alloc: 4.94MB
2014/07/20 12:52:06 7.819015362s for 1000 entities. Processed 34000, Mem alloc: 4.94MB
2014/07/20 12:52:14 7.608066645s for 1000 entities. Processed 35000, Mem alloc: 4.94MB
2014/07/20 12:52:22 7.711518627s for 1000 entities. Processed 36000, Mem alloc: 4.94MB
2014/07/20 12:52:29 7.658135451s for 1000 entities. Processed 37000, Mem alloc: 4.94MB
2014/07/20 12:52:38 8.251754200s for 1000 entities. Processed 38000, Mem alloc: 4.94MB
2014/07/20 12:52:46 8.296576343s for 1000 entities. Processed 39000, Mem alloc: 4.94MB
2014/07/20 12:52:54 8.230167582s for 1000 entities. Processed 40000, Mem alloc: 4.94MB
2014/07/20 12:53:02 7.796874592s for 1000 entities. Processed 41000, Mem alloc: 4.94MB
2014/07/20 12:53:10 7.954316905s for 1000 entities. Processed 42000, Mem alloc: 4.94MB
2014/07/20 12:53:18 7.814142008s for 1000 entities. Processed 43000, Mem alloc: 4.94MB
2014/07/20 12:53:26 7.863282043s for 1000 entities. Processed 44000, Mem alloc: 4.94MB
2014/07/20 12:53:33 7.695968977s for 1000 entities. Processed 45000, Mem alloc: 4.94MB
2014/07/20 12:53:40 6.960247879s for 1000 entities. Processed 46000, Mem alloc: 4.94MB
2014/07/20 12:53:47 6.975276363s for 1000 entities. Processed 47000, Mem alloc: 4.94MB
2014/07/20 12:53:55 7.855761010s for 1000 entities. Processed 48000, Mem alloc: 4.94MB
2014/07/20 12:54:00 Total articles: 48637
2014/07/20 12:54:00 Total errors: 0
2014/07/20 12:54:00 XML Parser took 0h 6m 371.544932s to run.
2014/07/20 12:54:00 XML Parser took 6m11.544931728s to run.
```

# PHP

see onixSQL.php

PHP needs > 12h of parsing. PHP 5.5 used.

To be fair that script should be rewritten to use the XmlReader()

# License

General Public License

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

Author
------

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

Made in Sydney, Australia :-)

If you consider a donation please contribute to: [http://www.seashepherd.org/](http://www.seashepherd.org/)
