<?php

ini_set('display_errors', 'On');
error_reporting(E_ALL);

$GLOBALS['THRIFT_ROOT'] = '/opt/app/thrift/lib/php';
$GLOBALS['SERVANT_ROOT'] = '../../servant/gen-php/fun/rpc';
require_once $GLOBALS['THRIFT_ROOT'].'/Thrift.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Base/TBase.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TTransportException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TProtocolException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TApplicationException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Protocol/TBinaryProtocol.php';
require_once $GLOBALS['THRIFT_ROOT'].'/StringFunc/TStringFunc.php';
require_once $GLOBALS['THRIFT_ROOT'].'/StringFunc/Core.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Factory/TStringFuncFactory.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Type/TType.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Type/TMessageType.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocket.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocketPool.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TBufferedTransport.php';
require_once $GLOBALS['SERVANT_ROOT'].'/FunServant.php';
require_once $GLOBALS['SERVANT_ROOT'].'/Types.php';

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
    $sock->setSendTimeout(1000);
    $sock->setRecvTimeout(2500);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123", 'reason' => 'call.init.121', 'host' => 'server1', 'ip' => '12.3.2.1'));

    // mysql select
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'SELECT * from UserInfo where uid>?', array(1));
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql update
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'UPDATE UserInfo set power=power+1 where uid=?', array(1));
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);

    // mysql merge blob column
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

    // mysql transation
    $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'BEGIN', NULL);
    $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'UPDATE UserInfo set power=power+1 where uid=?', array(1));
    $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'COMMIT', NULL);
    //$client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'ROLLBACK', NULL);
    echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
    print_r($rows);
    $transport->close();
} catch (TException $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
