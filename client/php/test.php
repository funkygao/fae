<?php

$GLOBALS['THRIFT_ROOT'] = '/opt/app/thrift/lib/php';
require_once $GLOBALS['THRIFT_ROOT'].'/Thrift.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Protocol/TBinaryProtocol.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocket.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/THttpClient.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TBufferedTransport.php';

$GEN_DIR = '../../servant/gen-php';  

require_once $GEN_DIR . '/fun/rpc/FunServant.php';  
require_once $GEN_DIR . '/fun/rpc/Types.php';  

try {
    $socket = new TSocket('localhost', 9001);
    $transport = new TBufferedTransport($socket, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new req_ctx(array('caller' => "me"));
    $return = $client->ping($ctx);
    echo $return;

    $transport->close();
} catch (TException $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
