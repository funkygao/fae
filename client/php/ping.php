<?php

require_once 'bootstrap.php';

try {
    $sock = new Thrift\Transport\TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(1);
    $sock->setSendTimeout(400000);
    $sock->setRecvTimeout(400000);
    $sock->setNumRetries(1);
    $transport = new Thrift\Transport\TBufferedTransport($sock, 4096, 4096);
    $protocol = new Thrift\Protocol\TBinaryProtocol($transport);

    // get our client
    $client = new fun\rpc\FunServantClient($protocol);
    $transport->open();

    $ctx = new fun\rpc\Context(array('rid' => hexdec(uniqid()), 'reason' => 'call.init.567', 'uid' => 11));

    for ($i=0; $i<2000; $i++) {
        echo 'ping:', $client->ping($ctx), "\n";
        echo 'noop:', $client->noop(21), "\n";
    }

    $transport->close();
} catch (Exception $ex) {
    print 'Something went wrong: ' . $ex->getMessage() . "\n";
}

