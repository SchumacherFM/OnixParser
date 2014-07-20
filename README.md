# OnixParser

This repo contains GoLang and PHP files to parse Onix xml files.

Can PHP beat GoLang? ;-)

# XML Data

[http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd](http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd)

#### Test data Onix Data Feed

You can download test data from [http://www.oup.com.au/help_and_advice/booksellers](http://www.oup.com.au/help_and_advice/booksellers)

The Complete File or The Incremental File

# Go

Inspired for learning by [http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)

Not all XML elements are matched because structure in DTD is unclear and we don't need all elements.

Output:

```
$ time ./OnixParserConcurrency --infile xmlFiles/oup_onix.xml --db test2 -v
OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm
This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html
2014/07/20 12:13:16 Dropped 19 existing tables
2014/07/20 12:13:18 1.879796552s Processed: 1000, child processes: 752, Mem alloc: 8.70MB
2014/07/20 12:13:21 2.759961177s Processed: 2000, child processes: 1255, Mem alloc: 10.39MB
2014/07/20 12:13:23 2.281196558s Processed: 3000, child processes: 1842, Mem alloc: 23.15MB
2014/07/20 12:13:26 2.2795532s Processed: 4000, child processes: 2440, Mem alloc: 23.02MB
2014/07/20 12:13:28 2.428083715s Processed: 5000, child processes: 3031, Mem alloc: 37.22MB
2014/07/20 12:13:30 2.015324302s Processed: 6000, child processes: 3709, Mem alloc: 37.75MB
2014/07/20 12:13:32 2.020671416s Processed: 7000, child processes: 4426, Mem alloc: 43.50MB
2014/07/20 12:13:34 2.20758561s Processed: 8000, child processes: 4995, Mem alloc: 35.19MB
2014/07/20 12:13:36 1.94296616s Processed: 9000, child processes: 5681, Mem alloc: 44.15MB
2014/07/20 12:13:38 2.084476139s Processed: 10000, child processes: 6310, Mem alloc: 76.65MB
2014/07/20 12:13:40 1.867831049s Processed: 11000, child processes: 6999, Mem alloc: 65.31MB
2014/07/20 12:13:42 1.837092633s Processed: 12000, child processes: 7804, Mem alloc: 50.57MB
2014/07/20 12:13:44 1.964309595s Processed: 13000, child processes: 8584, Mem alloc: 81.21MB
2014/07/20 12:13:46 2.108973939s Processed: 14000, child processes: 9358, Mem alloc: 58.65MB
2014/07/20 12:13:50 3.455177739s Processed: 15000, child processes: 10022, Mem alloc: 89.67MB
2014/07/20 12:13:52 2.401940269s Processed: 16000, child processes: 10738, Mem alloc: 108.28MB
2014/07/20 12:13:54 1.587906304s Processed: 17000, child processes: 11544, Mem alloc: 122.63MB
2014/07/20 12:13:56 2.072638953s Processed: 18000, child processes: 12281, Mem alloc: 133.29MB
2014/07/20 12:13:58 2.074804128s Processed: 19000, child processes: 12990, Mem alloc: 143.84MB
2014/07/20 12:14:00 2.159436142s Processed: 20000, child processes: 13656, Mem alloc: 148.84MB
2014/07/20 12:14:02 2.414585834s Processed: 21000, child processes: 14243, Mem alloc: 153.96MB
2014/07/20 12:14:05 2.120901052s Processed: 22000, child processes: 14847, Mem alloc: 154.85MB
2014/07/20 12:14:07 2.24752107s Processed: 23000, child processes: 15346, Mem alloc: 150.62MB
2014/07/20 12:14:09 2.368137206s Processed: 24000, child processes: 15941, Mem alloc: 138.14MB
2014/07/20 12:14:11 2.102820803s Processed: 25000, child processes: 16496, Mem alloc: 109.79MB
2014/07/20 12:14:14 2.330307509s Processed: 26000, child processes: 16838, Mem alloc: 189.57MB
2014/07/20 12:14:16 2.349806224s Processed: 27000, child processes: 17109, Mem alloc: 169.57MB
2014/07/20 12:14:18 2.162176411s Processed: 28000, child processes: 17694, Mem alloc: 138.28MB
2014/07/20 12:14:20 2.282641359s Processed: 29000, child processes: 18290, Mem alloc: 214.77MB
2014/07/20 12:14:23 2.192836431s Processed: 30000, child processes: 19037, Mem alloc: 182.17MB
2014/07/20 12:14:25 2.368323714s Processed: 31000, child processes: 19608, Mem alloc: 151.61MB
2014/07/20 12:14:27 1.633814548s Processed: 32000, child processes: 20345, Mem alloc: 221.78MB
2014/07/20 12:14:29 1.987422489s Processed: 33000, child processes: 20988, Mem alloc: 180.15MB
2014/07/20 12:14:31 1.99859052s Processed: 34000, child processes: 21626, Mem alloc: 260.88MB
2014/07/20 12:14:33 1.955077052s Processed: 35000, child processes: 22327, Mem alloc: 211.58MB
2014/07/20 12:14:35 2.071890766s Processed: 36000, child processes: 23085, Mem alloc: 230.58MB
2014/07/20 12:14:37 1.815911816s Processed: 37000, child processes: 23829, Mem alloc: 232.24MB
2014/07/20 12:14:39 2.021806767s Processed: 38000, child processes: 24644, Mem alloc: 221.69MB
2014/07/20 12:14:40 1.707361916s Processed: 39000, child processes: 25467, Mem alloc: 241.95MB
2014/07/20 12:14:42 1.867942278s Processed: 40000, child processes: 26242, Mem alloc: 319.55MB
2014/07/20 12:14:44 2.14585482s Processed: 41000, child processes: 27045, Mem alloc: 246.45MB
2014/07/20 12:14:46 1.982444273s Processed: 42000, child processes: 27705, Mem alloc: 326.24MB
2014/07/20 12:14:48 1.969034462s Processed: 43000, child processes: 28480, Mem alloc: 240.94MB
2014/07/20 12:14:50 1.927976522s Processed: 44000, child processes: 29073, Mem alloc: 319.64MB
2014/07/20 12:14:52 2.076284518s Processed: 45000, child processes: 29494, Mem alloc: 222.37MB
2014/07/20 12:14:54 1.832563383s Processed: 46000, child processes: 30030, Mem alloc: 294.21MB
2014/07/20 12:14:55 1.380362894s Processed: 47000, child processes: 30623, Mem alloc: 359.07MB
2014/07/20 12:14:57 1.843187116s Processed: 48000, child processes: 31224, Mem alloc: 240.56MB
2014/07/20 12:15:08 29499 child processes remaining ... 2014-07-20 12:15:08.908702699 +1000 EST
2014/07/20 12:15:18 26725 child processes remaining ... 2014-07-20 12:15:18.908334787 +1000 EST
2014/07/20 12:15:28 23132 child processes remaining ... 2014-07-20 12:15:28.908345311 +1000 EST
2014/07/20 12:15:38 19391 child processes remaining ... 2014-07-20 12:15:38.908344404 +1000 EST
2014/07/20 12:15:48 13802 child processes remaining ... 2014-07-20 12:15:48.908354984 +1000 EST
2014/07/20 12:15:58 6 child processes remaining ... 2014-07-20 12:15:58.909032752 +1000 EST
2014/07/20 12:15:58 Total articles: 48637
2014/07/20 12:15:58 Total errors: 0
2014/07/20 12:15:58 XML Parser took 0h 2m 162.285939s to run.
2014/07/20 12:15:58 XML Parser took 2m42.285938833s to run.

real	2m42.345s
user	1m40.032s
sys	    0m58.290s
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
