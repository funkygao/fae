<?php

require_once 'bootstrap.php';

try {
    $sock = new Thrift\Transport\TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(1);
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new Thrift\Transport\TBufferedTransport($sock, 1024, 1024);
    $protocol = new Thrift\Protocol\TBinaryProtocol($transport);

    // get our client
    $client = new fun\rpc\FunServantClient($protocol);
    $transport->open();

    $ctx = new fun\rpc\Context(array('rid' => hexdec(uniqid()), 
        'reason' => 'test.couchbase'));

    // couchbase get/set
    $bucket = 'default';
    for ($i=0; $i<10000; $i++) {
        $ok = $client->cb_set($ctx, $bucket, 'key1', 'value1', 0);

        $value = $client->cb_get($ctx, $bucket, 'key1');
        var_dump($value);
    }

    $transport->close();
} catch (Exception $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
