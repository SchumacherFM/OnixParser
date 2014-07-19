# OnixParser

This repo contains GoLang and PHP files to parse Onix xml files.

Can PHP beat GoLang? ;-)

# XML Data

the original availability onix xml file is 3GB huge.

[http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd](http://www.editeur.org/onix/2.1/02/reference/onix-international.dtd)

# Go

Inspired for learning by [http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)

Average import time for 1000 Products: ~6 sec

Not all XML elements are matched because structure in DTD is unclear and we don't need all elements.

# PHP

see onixSQL.php

PHP needs > 12h of parsing. PHP 5.5 used.

# License

General Public License

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)
