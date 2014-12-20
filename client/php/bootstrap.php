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

