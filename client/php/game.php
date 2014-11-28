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
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => "123", 'reason' => 'call.init', 'host' => 'server1', 'ip' => '12.3.2.1'));

    // game get unique name with len 3
    for ($i = 0; $i < 100; $i ++) {
        echo $client->gm_name3($ctx), "\n";
    }

    $transport->close();
} catch (TException $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
