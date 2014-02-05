<?php

// must include mongo extension before use this script

$data = array(
    'uid' => 12122,
    'abtype' => array(
        'tutorial' => 'a',
        'payment' => 'b',
    ),
    'gendar' => 'F',
);

$pack = bson_encode($data);
echo $pack, "\n";
print_r(bson_decode($pack));

