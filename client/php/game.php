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
    $sock = new TSocketPool(array('localhost'), array(9011));
    $sock->setDebug(1);
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123", 'reason' => 'call.init.567', 'host' => 'server1', 'ip' => '12.3.2.1'));

    // game get unique name with len 3
    for ($i = 0; $i < 5; $i ++) {
        echo $client->gm_name3($ctx), "\n";
    }

    $ok = $client->zk_create($ctx, "/maintain/global", "");
    var_dump($ok);
    $nodes = $client->zk_children($ctx, "/maintain");
    print_r($nodes);
    $ok = $client->zk_del($ctx, "/maintain/global");
    var_dump($ok);

    $transport->close();
} catch (Exception $ex) {
    print 'Something went wrong: ' . $ex->getMessage() . "\n";
}

