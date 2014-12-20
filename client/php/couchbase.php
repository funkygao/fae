<?php

require_once 'bootstrap.php';

try {
    $sock = new TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(1);
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123nfa", 'reason' => 'test.couchbase', 'host' => 'phptest', 'ip' => '12.3.2.1'));

    // couchbase get/set
    $bucket = 'default';
    for ($i=0; $i<10000; $i++) {
        $ok = $client->cb_set($ctx, $bucket, 'key1', 'value1', 0);

        $value = $client->cb_get($ctx, $bucket, 'key1');
        var_dump($value);
    }

    $transport->close();
} catch (TException $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
