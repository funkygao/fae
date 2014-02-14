<?php
/**
 * json_encode vs serialize performance
 */


ini_set('display_errors', 1);
error_reporting(E_ALL);

//  Make a bit, honkin test array
//  You may need to adjust this depth to avoid memory limit errors
$testArray = fillArray(0, 5);

//  Time json encoding
$start = microtime(true);
json_encode($testArray);
$jsonTime = microtime(true) - $start;
echo "JSON encoded in $jsonTime seconds\n";

//  Time serialization
$start = microtime(true);
serialize($testArray);
$serializeTime = microtime(true) - $start;
echo "PHP serialized in $serializeTime seconds\n";

//  Compare them
if ( $jsonTime < $serializeTime ) {
    echo "json_encode() was roughly " . 
        number_format( ($serializeTime / $jsonTime - 1 ) * 100, 2 ) . "% faster than serialize()\n";
} else if ( $serializeTime < $jsonTime ) {
    echo "serialize() was roughly " . 
        number_format( ($jsonTime / $serializeTime - 1 ) * 100, 2 ) . "% faster than json_encode()\n";
} else {
    echo 'Unpossible!\n';
}

function fillArray( $depth, $max ) {
    static $seed;
    if ( is_null( $seed ) ) {
        $seed = array( 'a', 2, 'c', 4, 'e', 6, 'g', 8, 'i', 10 );
    }
    if ( $depth < $max ) {
        $node = array();
        foreach ( $seed as $key ) {
            $node[$key] = fillArray( $depth + 1, $max );
        }
        return $node;
    }

    return 'empty';
}

