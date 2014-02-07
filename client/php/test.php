<?php

$GLOBALS['THRIFT_ROOT'] = '/opt/app/thrift/lib/php';
$GLOBALS['SERVANT_ROOT'] = '../../servant/gen-php/fun/rpc';
require_once $GLOBALS['THRIFT_ROOT'].'/Thrift.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Exception/TException.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Protocol/TBinaryProtocol.php';
require_once $GLOBALS['THRIFT_ROOT'].'/StringFunc/TStringFunc.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Factory/TStringFuncFactory.php';
require_once $GLOBALS['THRIFT_ROOT'].'/StringFunc/Core.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Type/TType.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Type/TMessageType.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocket.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TSocketPool.php';
require_once $GLOBALS['THRIFT_ROOT'].'/Transport/TBufferedTransport.php';
require_once $GLOBALS['SERVANT_ROOT'].'/FunServant.php';
require_once $GLOBALS['SERVANT_ROOT'].'/Types.php';

try {
    $sock = new Thrift\Transport\TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(0);
    $sock->setSendTimeout(1000);
    $sock->setRecvTimeout(2500);
    $sock->setNumRetries(1);
    $transport = new Thrift\Transport\TBufferedTransport($sock, 1024, 1024);
    $protocol = new Thrift\Protocol\TBinaryProtocol($transport);

    // get our client
    $client = new fun\rpc\FunServantClient($protocol);
    $transport->open();

    $ctx = new fun\rpc\Context(array('caller' => "me"));
    $return = $client->ping($ctx);
    echo $return;

    $transport->close();
} catch (TException $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}
