<?php
##########################################################################
##########################################################################
##                                                                      ##
##  This script processes ONIX xml files and inserts them into a        ##
##  database. This database does not need to have any pre-existing      ##
##  tables or collums, these will be automatically created by the       ##
##  script.                                                             ##
##                                                                      ##
##                                                                      ##
##                                                                      ##
##  AFTER RUNNING THIS SCRIPT:                                          ##
##  update table collums etc. to match content type.                    ##
##  Also you might want to ad some primary keys and indexes afterwards  ##
##                                                                      ##
##  Author: Jonathan van Bochove                                        ##
##  Author url: www.johannes-multimedia.nl                              ##
##  Author e-mail: webmaster@johannes-multimedia.nl                     ##
##  Licence: Copyright (c) 2011 Johannes Multimedia                     ##
##  Released under the GNU General Public License                       ##
##  Version 1.2.2 (2013-02-18)                                          ##
##  Version 1.3 (2014-06-26) bugfixes, mysqli (CS)                      ##
##                                                                      ##
##                                                                      ##
##  If you make any alterations to this script to make it more use-     ##
##  full, faster or more efficient, please send a copy of the updated   ##
##  script to the author, and mention what was updated/changed.         ##
##                                                                      ##
##                                                                      ##
##########################################################################
##                                                                      ##
##  edit settings below:                                                ##
##                                                                      ##
##########################################################################
##########################################################################

$mem        = 1000000; // Onix chunk size in bytes (script won't process more then this at once)
$file       = "demo-availability.xml"; // Location of onix file
$dbHost     = "localhost"; // mysql host
$dbUser     = "root"; //mysql username
$dbPassword = "325639de2967f22b920536f49355825ce07e2a46"; // mysql user password
$db         = "test"; // mysql database name
$prefix     = "onix_"; // table prefix

##########################################################################
##########################################################################
##                                                                      ##
##                           end of settings                            ##
##                                                                      ##
##########################################################################
##########################################################################
mysqli_report(MYSQLI_REPORT_STRICT);
try {
    $mysqli = new mysqli($dbHost, $dbUser, $dbPassword, $db);
} catch (Exception $e) {
    die($e->getMessage() . PHP_EOL);
}

function ti()
{ // function to calculate time used
    $t = microtime();
    $t = explode(' ', $t);
    return $t[1] + $t[0];
}

$argvStart = (int)(isset($argv[1]) ? $argv[1] : 0);
$argvSt    = (float)(isset($argv[2]) ? $argv[2] : 0);
$argvTotal = (float)(isset($argv[3]) ? $argv[3] : 0);

$start = max($argvStart, 0); // set startpoint of xml file (must be an integer greater then 0)
$st    = $start === 0 ? ti() : $argvSt; // remember when we started with the first chunk of data to show total processing time
$total = max($argvTotal, 0); // remember the total number of records we processed from the start of the first chunk
$size  = filesize($file);

$end = min(($size - $start), $mem); // if xml file smaller then chunksize, then don't try and do too much

if ($start < $size) { // are we not already done?
    $p       = file_get_contents($file, null, null, $start, $end); // load the chunk of xml into memory
    $pos     = strripos($p, '</Product>') + 10; // find the initial end of the last record of this chunk of data
    $deleted = strlen($p) - $pos; // help to figure out where to start processing the next chunk of data
    $p       = preg_replace("/(.*?)<([Pp])roduct(.*?)>(.*)/s", "<\\2roduct>\\4", $p); // strip the "useless" header and stuff
    $p       = preg_replace("!<br />!", "&#60;br /&#62;", $p); //turning possible <br /> html into its special chars equivalent
    $p       = preg_replace('!<([^ ]*?) ([^=]*?)="([^"]*?)">!', '<\\2>\\3</\\2><\\1>', $p); //turning tag values into their own tags
    $pos     = strripos($p, '</Product>') + 10; // find the end of the last record of this chunk of data, after modifications from above
    $product = '';

    if ($pos > 10) { //are there more products to process?
        $products = simplexml_load_string("<xml>" . substr($p, 0, $pos) . "</xml>"); // do the magic, turn the xml into an xml object that we can process
        unset($p); // clear the memory of the xml string
        if (is_object($products)) {
            $total = $total + count($products); // how many records to process?
        }

        // Fetch existing tables and columns
        $tbl  = array();
        $tbls = $mysqli->query("SHOW TABLES like '" . $prefix . "%'");
        while ($temp = $tbls->fetch_array()) {
            $tbl[strtolower($temp[0])] = array();
        }
        foreach ($tbl as $key => $value) {
            $columns = $mysqli->query("SHOW COLUMNS FROM `" . $mysqli->real_escape_string($key) . "`");
            while ($temp = $columns->fetch_array()) {
                $tbl[strtolower($key)][$temp['Field']] = " ";
            }
        }

        // if it does not exist, create the first table
        if (!isset($tbl[$prefix . 'product'])) {
            $mysqli->query("CREATE TABLE IF NOT EXISTS `" . $mysqli->real_escape_string($prefix . "product") . "` (`id` varchar(15) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8");
            $tbl[$prefix . 'product']['id'] = " ";
        }
        // loop through the chunk of xml to process
        foreach ($products as $produc) {
            //check that all used tables and collumns exist, and if not create and update these tables
            if (isset($produc->productidentifier)) {
                foreach ($produc->productidentifier as $value) {
                    if ($value->b221 == '03' || $value->b221 == '15') $id = $mysqli->real_escape_string($value->b244);
                }
            } else {
                foreach ($produc->ProductIdentifier as $value) {
                    if ($value->ProductIDType == '03' || $value->ProductIDType == '15') $id = $mysqli->real_escape_string($value->IDValue);
                }
            }
            $varup = null;
            foreach ($produc as $key => $value) { //loop trough everything building database and writing insert queue
                $vars = get_object_vars($value);
                if (is_array($vars) && count($vars) > 0) {
                    $i     = ($key === $varup ? ($i + 1) : 0); //count the number of instances of a certain tag
                    $varup = $key;
                    $key   = strtolower($prefix . $key); //table names must be lowercase, with prefix prepended
                    if (!isset($tbl[$key])) { // create missing tables
                        $mysqli->query("CREATE TABLE IF NOT EXISTS `" . $mysqli->real_escape_string($key) . "` (`id` varchar(15) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8");
                        $tbl[$key] = array('id' => 'varchar(15)');
                    }
                    $varup2 = null;
                    $j      = 0;
                    foreach ($value as $key2 => $value2) {
                        $vars2 = get_object_vars($value2);
                        if (is_array($vars2) && count($vars2) > 0) {
                            $j      = ($varup2 === $key ? (int)($j + 1) : $j); //count the number of instances of a certain lvl2 tag
                            $varup2 = $key;
                            foreach ($value2 as $key3 => $value3) {
                                if (!isset($tbl[$key][$key3])) { //add missing columns to tables
                                    $mysqli->query("ALTER TABLE `" . $mysqli->real_escape_string($key) . "` ADD `" . $mysqli->real_escape_string($key3) . "` longtext");
                                    $tbl[$key][$key3] = 'longtext';
                                }
                                ${$key}[$id][$j][$key3] = (string)$value3; //write queue for lvl 3 tags, I don't know of any lvl 4 tags so I stop seaarching here, please correct me if I'm wrong
                            }
                        } else {
                            if (!isset($tbl[$key][$key2])) { //add missing columns to tables
                                $mysqli->query("ALTER TABLE `" . $mysqli->real_escape_string($key) . "` ADD `" . $mysqli->real_escape_string($key2) . "` longtext");
                                $tbl[$key][$key2] = 'longtext';
                            }
                            ${$key}[$id][$i][$key2] = (string)$value2; //write queue for lvl 2 tags
                        }
                    }
                } else {
                    if (!isset($tbl[$prefix . 'product'][$key])) { //update primary table for missing columns
                        $mysqli->query("ALTER TABLE " . $prefix . "product ADD " . $mysqli->real_escape_string($key) . " VARCHAR(128)");
                        $tbl[$prefix . 'product'][$key] = 'varchar(128)';
                    }
                    $i                                   = 0; //I don't know of any recurring first lvl tags, so $i is always 0, correct me if I'm wrong
                    ${$prefix . "product"}[$id][0][$key] = (string)$value; //write queue for first lvl tags
                }
            }
        }
        foreach ($tbl as $table => $array1) { // check if we can save some inserts by merging records
            if (isset(${$table}) && is_array(${$table}) && count(${$table}) > 0) {
                foreach (${$table} as $key => $array) {
                    if (count($array) > 0) {
                        sort(${$table}[$key]);
                        sort($array);
                        for ($a = (count($array) - 1); $a > 0; $a--) {
                            $test = array_merge($array[$a], $array[($a - 1)]);
                            if ((count($array[$a]) + count($array[($a - 1)])) == count($test)) {
                                ${$table}[$key][($a - 1)] = $test;
                                $array[($a - 1)]          = $test;
                                unset(${$table}[$key][$a], $array[$a]);
                            }
                        }
                    }
                }
            }
        }
        // insert each array of data into its own table
        foreach ($tbl as $table => $array) {
            $query = "REPLACE INTO `" . $mysqli->real_escape_string($table) . "` (";
            foreach ($array as $key => $useless) {
                $query .= "`" . $mysqli->real_escape_string($key) . "`, ";
            }
            $query = substr($query, 0, -2) . ") VALUES ";
            $rows  = '';
            if (isset(${$table}) && is_array(${$table}) && count(${$table}) > 0) {
                foreach (${$table} as $key => $value) {
                    foreach ($value as $key2 => $value2) {
                        # $key = isbn
                        # $value2 = array with data to be inserted
                        $rows .= '(';
                        foreach ($array as $k => $v) {
                            $rows .= "'" . ($k == 'id'
                                    ? $mysqli->real_escape_string($key)
                                    : (isset($value2[$k])
                                        ? $mysqli->real_escape_string(utf8_decode($value2[$k]))
                                        : '')) . "', ";
                        }
                        $rows = substr($rows, 0, -2) . "), ";
                    }
                }
            }
            $rows         = substr($rows, 0, -2);
            $insertResult = $mysqli->query($query . $rows);
            if (false === $insertResult) {
                echo "FAILED: $query$rows\n";
            }
        }
    }
    if (($end + $start - ($deleted + 1)) < $size) {
        // continue with the next chunk of xml
        $command = 'nohup php ' . $_SERVER['PHP_SELF'] . " " . ($end + $start - ($deleted + 1)) . " " . $st . " " . $total . ' >> onixLog.txt &';
        echo "Starting: $command\n\n";
        exec($command);
    } else { // finished and show total number of inserted records and processing time
        echo date("Y-m-d H:i:s") . " records: " . $total . " time: " . number_format((ti() - $st), 2, '.', ',') . " seconds";
    }
    $mysqli->close();
}
