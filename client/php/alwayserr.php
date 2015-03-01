<?php

require_once 'bootstrap.php';

try {
    $sock = new Thrift\Transport\TSocketPool(array('localhost', 'localhost'), array(9001, 9011));
    $sock->setDebug(1);
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new Thrift\Transport\TBufferedTransport($sock, 1024, 1024);
    $protocol = new Thrift\Protocol\TBinaryProtocol($transport);

    // get our client
    $client = new fun\rpc\FunServantClient($protocol);
    $transport->open();

    $ctx = new fun\rpc\Context(array('rid' => "123", 'reason' => 'call.init.121'));

    // mysql select multiple rows
    for ($i = 0; $i < 10; $i++) {
        $uid = 1;
        $invalidUid = 0; // non-exist
        try {
            echo PHP_EOL, $i+1, PHP_EOL;

            $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', $uid, 'SELECT uid from UserInfo where uid=?', array($uid), '');
            echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;

            $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', $invalidUid, 'SELECT uid from UserInfo where uid=?', array($invalidUid), '');
            echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;

        } catch (Thrift\Exception\TApplicationException $ex) {
            echo $ex->getMessage(), "\n";
        }
    }

    $transport->close();
} catch (Exception $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}

