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
    $sock = new TSocketPool(array('localhost', 'localhost'), array(9001, 9011));
    $sock->setDebug(1);
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => hexdec(uniqid()), 'reason' => 'call.init.121'));

    // mysql select multiple rows
    echo "\nDEMO SELECT\n";
    echo "===============================\n";
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'SELECT * from UserInfo where uid>?', array(1), '');
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql query cache
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'SELECT * from UserInfo where uid=?', array(1), 'UserInfo:1');
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'SELECT * from UserInfo where uid=?', array(1), 'UserInfo:1');
    print_r($rows);

    // mysql update
    echo "\nDEMO UPDATE\n";
    echo "===============================\n";
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'UPDATE UserInfo set power=power+1 where uid=?', array(1), 'UserInfo:1');
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql transation
    echo "\nDEMO transtaion\n";
    echo "===============================\n";
    $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'BEGIN', NULL, '');
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'UPDATE UserInfo set power=power+1 where uid=?', array(1), 'UserInfo:1');
    $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'COMMIT', NULL, '');
    //$client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'ROLLBACK', NULL, '');
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql bulk exec
    echo "\nDEMO bulk exec\n";
    echo "===============================\n";
    $client->my_bulk_exec($ctx, 
        array('UserShard', 'AllianceShard'),
        array('UserInfo', 'Alliance'),
        array(1, 3),
        array(
            'UPDATE UserInfo set power=? WHERE uid=?',
            'UPDATE Alliance set power=? WHERE alliance_id=?',
        ),
        array(
            array(158, 1),
            array(1508, 3),
        ),
        array(
            '', 
            '',
        )
    );

    // mysql query shards
    echo "\nDEMO query shards\n";
    echo "===============================\n";
    $rows = $client->my_query_shards($ctx, 'UserShard', 'UserInfo', 'SELECT chat_channel FROM UserInfo WHERE uid>?', array(1));
    print_r($rows);

    // mysql merge blob column
    echo "\nDEMO MERGE\n";
    echo "===============================\n";
    $merged = $client->my_merge($ctx, 'AllianceShard', 'Rally', 1, 'alliance_id=51 and uid=50', 
        'Rally:' . json_encode(array(
            'alliance_id' => 51,
            'uid' => 50,
        )),
        'slots_info', 
        json_encode(
            array(
                'info' => array( 
                    "88" => time(),
                )
            )));
    print_r($merged);
    print_r(json_decode($merged->newVal, TRUE));

    $transport->close();
} catch (Exception $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}

