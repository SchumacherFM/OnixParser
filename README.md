# OnixParser Concurrent execution

This repo contains GoLang code for concurrent execution.

Binary is build for OSX 10.9 Darwin.

# XML Data

[http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd](http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd)

#### Test data Onix Data Feed

You can download test data from [http://www.oup.com.au/help_and_advice/booksellers](http://www.oup.com.au/help_and_advice/booksellers)

The Complete File or The Incremental File

# Go

Inspired for learning by [http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)

Not all XML elements are matched because structure in DTD is unclear and we don't need all elements.

XML Filesize 278MB Output:

```
$ ./OnixParser --infile xmlFiles/oup_onix.xml --db test2 -v
OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm
This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html
2014/07/20 12:58:28 Dropped 15 existing tables
2014/07/20 12:58:30 1.696429991s Processed: 1000, child processes: 783, Mem alloc: 22.70MB
2014/07/20 12:58:32 2.000549676s Processed: 2000, child processes: 1425, Mem alloc: 37.40MB
2014/07/20 12:58:34 2.032047269s Processed: 3000, child processes: 2014, Mem alloc: 50.28MB
2014/07/20 12:58:36 2.126721153s Processed: 4000, child processes: 2719, Mem alloc: 64.66MB
2014/07/20 12:58:39 2.477956106s Processed: 5000, child processes: 3347, Mem alloc: 78.42MB
2014/07/20 12:58:41 2.366202751s Processed: 6000, child processes: 3912, Mem alloc: 92.49MB
2014/07/20 12:58:44 2.511615940s Processed: 7000, child processes: 4562, Mem alloc: 103.99MB
2014/07/20 12:58:46 2.360495760s Processed: 8000, child processes: 5159, Mem alloc: 119.29MB
2014/07/20 12:58:48 1.876899327s Processed: 9000, child processes: 5753, Mem alloc: 131.36MB
2014/07/20 12:58:49 1.678509629s Processed: 10000, child processes: 6450, Mem alloc: 145.37MB
2014/07/20 12:58:52 2.189865604s Processed: 11000, child processes: 7005, Mem alloc: 154.37MB
2014/07/20 12:58:54 2.063110983s Processed: 12000, child processes: 7718, Mem alloc: 165.06MB
2014/07/20 12:58:55 1.634053347s Processed: 13000, child processes: 8485, Mem alloc: 176.38MB
2014/07/20 12:58:57 2.095966550s Processed: 14000, child processes: 9218, Mem alloc: 192.37MB
2014/07/20 12:58:59 1.852198041s Processed: 15000, child processes: 9977, Mem alloc: 204.26MB
2014/07/20 12:59:01 1.796805218s Processed: 16000, child processes: 10743, Mem alloc: 219.26MB
2014/07/20 12:59:03 2.067620401s Processed: 17000, child processes: 11492, Mem alloc: 237.65MB
2014/07/20 12:59:05 1.619303423s Processed: 18000, child processes: 12258, Mem alloc: 253.60MB
2014/07/20 12:59:07 1.955929824s Processed: 19000, child processes: 13010, Mem alloc: 268.31MB
2014/07/20 12:59:09 2.074955765s Processed: 20000, child processes: 13754, Mem alloc: 284.82MB
2014/07/20 12:59:11 2.295056187s Processed: 21000, child processes: 14438, Mem alloc: 301.32MB
2014/07/20 12:59:13 2.115669885s Processed: 22000, child processes: 15059, Mem alloc: 310.08MB
2014/07/20 12:59:15 2.094373829s Processed: 23000, child processes: 15624, Mem alloc: 323.09MB
2014/07/20 12:59:17 2.098578482s Processed: 24000, child processes: 16135, Mem alloc: 339.16MB
2014/07/20 12:59:19 1.438987846s Processed: 25000, child processes: 16920, Mem alloc: 360.17MB
2014/07/20 12:59:21 1.704874825s Processed: 26000, child processes: 17446, Mem alloc: 365.61MB
2014/07/20 12:59:23 2.187127534s Processed: 27000, child processes: 17750, Mem alloc: 374.87MB
2014/07/20 12:59:25 1.935661751s Processed: 28000, child processes: 18212, Mem alloc: 391.25MB
2014/07/20 12:59:27 2.228236543s Processed: 29000, child processes: 18744, Mem alloc: 399.50MB
2014/07/20 12:59:29 1.807835524s Processed: 30000, child processes: 19312, Mem alloc: 410.88MB
2014/07/20 12:59:31 2.057599998s Processed: 31000, child processes: 20096, Mem alloc: 429.52MB
2014/07/20 12:59:32 1.602475634s Processed: 32000, child processes: 20793, Mem alloc: 435.14MB
2014/07/20 12:59:34 1.937873196s Processed: 33000, child processes: 21550, Mem alloc: 456.59MB
2014/07/20 12:59:38 3.187732317s Processed: 34000, child processes: 22116, Mem alloc: 471.67MB
2014/07/20 12:59:40 2.207788615s Processed: 35000, child processes: 22873, Mem alloc: 485.80MB
2014/07/20 12:59:42 1.984706182s Processed: 36000, child processes: 23576, Mem alloc: 510.57MB
2014/07/20 12:59:44 1.998194673s Processed: 37000, child processes: 24294, Mem alloc: 517.38MB
2014/07/20 12:59:46 1.909828674s Processed: 38000, child processes: 25018, Mem alloc: 543.34MB
2014/07/20 12:59:47 1.658589319s Processed: 39000, child processes: 25829, Mem alloc: 551.84MB
2014/07/20 12:59:49 2.043943403s Processed: 40000, child processes: 26642, Mem alloc: 562.34MB
2014/07/20 12:59:52 2.175651195s Processed: 41000, child processes: 27389, Mem alloc: 588.97MB
2014/07/20 12:59:54 2.028642737s Processed: 42000, child processes: 28149, Mem alloc: 599.48MB
2014/07/20 12:59:56 2.075818822s Processed: 43000, child processes: 28831, Mem alloc: 625.18MB
2014/07/20 12:59:58 2.232703332s Processed: 44000, child processes: 29507, Mem alloc: 630.81MB
2014/07/20 13:00:00 2.384482530s Processed: 45000, child processes: 30046, Mem alloc: 663.90MB
2014/07/20 13:00:02 1.763376000s Processed: 46000, child processes: 30616, Mem alloc: 668.52MB
2014/07/20 13:00:04 1.577886029s Processed: 47000, child processes: 31203, Mem alloc: 672.77MB
2014/07/20 13:00:06 2.179890588s Processed: 48000, child processes: 31700, Mem alloc: 705.86MB
2014/07/20 13:00:17 29936 child processes remaining ... 2014-07-20 13:00:17.386865822 +1000 EST
2014/07/20 13:00:27 27717 child processes remaining ... 2014-07-20 13:00:27.386898857 +1000 EST
2014/07/20 13:00:37 24315 child processes remaining ... 2014-07-20 13:00:37.386867606 +1000 EST
2014/07/20 13:00:47 20820 child processes remaining ... 2014-07-20 13:00:47.386878769 +1000 EST
2014/07/20 13:00:57 16592 child processes remaining ... 2014-07-20 13:00:57.38693394 +1000 EST
2014/07/20 13:01:07 7208  child processes remaining ... 2014-07-20 13:01:07.386874033 +1000 EST
2014/07/20 13:01:17 6     child processes remaining ... 2014-07-20 13:01:17.387629411 +1000 EST
2014/07/20 13:01:17 Total articles: 48637
2014/07/20 13:01:17 Total errors: 0
2014/07/20 13:01:17 XML Parser took 0h 2m 169.110496s to run.
2014/07/20 13:01:17 XML Parser took 2m49.110496443s to run.
```

There are severeal options on the command line:

```
$ go run OnixParser.go -h
OnixParser Copyright (C) 2014 Cyrill AT Schumacher dot fm
This program comes with ABSOLUTELY NO WARRANTY; License: http://www.gnu.org/copyleft/gpl.html
Usage of OnixParser:
  -db="test": MySQL db name
  -host="127.0.0.1": MySQL host name
  -infile="": Input file path
  -mla=6.5: Max Load Average, float value. Recommended > 6, if <= 3 then disabled
  -moc=20: Max MySQL open connections
  -pass="test": MySQL password
  -tablePrefix="gonix_": Table name prefix
  -user="test": MySQL user name
  -v=false: Increase verbosity
exit status 2
```

Checking for the Max Load Average means if the Load AVG will be above that value the programm will stop and wait until the Load AVG will fall below that threshold.

# License

General Public License

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

Author
------

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

Made in Sydney, Australia :-)

If you consider a donation please contribute to: [http://www.seashepherd.org/](http://www.seashepherd.org/)
