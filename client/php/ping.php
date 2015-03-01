<?php

require_once 'bootstrap.php';

use Thrift\Transport\TSocketPool;
use Thrift\Transport\TBufferedTransport;
use Thrift\Protocol\TBinaryProtocol;
use Thrift\Exception\TTransportException;
use Thrift\Exception\TProtocolException;
use fun\rpc\FunServantClient;
use fun\rpc\Context;
use fun\rpc\TCacheMissed;
use fun\rpc\TMongoMissed;
use fun\rpc\TMemcacheData;

try {
    $sock = new TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(1);
    $sock->setSendTimeout(400000);
    $sock->setRecvTimeout(400000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 4096, 4096);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123", 'reason' => 'call.init.567', 'uid' => 11));

    echo $client->ping($ctx), "\n";

    $client->gm_latency($ctx, 19, 21);
    var_dump($client->gm_presence($ctx, array(11, 14)));

    $transport->close();
} catch (Exception $ex) {
    print 'Something went wrong: ' . $ex->getMessage() . "\n";
}

