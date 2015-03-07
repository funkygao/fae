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
    $sock = new TSocketPool(array('localhost'), array(9001));
    $sock->setDebug(1);
    $sock->setSendTimeout(400000);
    $sock->setRecvTimeout(400000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 4096, 4096);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('rid' => hexdec(uniqid()), 'reason' => 'call.init.567'));

    $client->gm_latency($ctx, 19, 21);
    var_dump($client->gm_presence($ctx, array(11, 14)));

    var_dump($client->gm_reserve($ctx, 'u', 'funky1', 'funky'));
    var_dump($client->gm_reserve($ctx, 'u', 'funky1', 'funky'));
    var_dump($client->gm_reserve($ctx, 'u', 'funky', 'funky1'));
    var_dump($client->gm_reserve($ctx, 'u', 'funky', 'funky1'));

    for ($i=0; $i<100; $i++) {
        $k = $client->gm_register($ctx, 'k');
        echo "reg with k: $k \n";
    }

    // redis
    $r = $client->rd_call($ctx, 'get', 'default', array('_not_existent_key'));
    var_dump($r);
    $r = $client->rd_call($ctx, 'set', 'default', array('the key', 'hello world!',));
    var_dump($r);
    $r = $client->rd_call($ctx, 'get', 'default', array('the key'));
    var_dump($r);
    $r = $client->rd_call($ctx, 'del', 'default', array('the key'));
    $client->rd_call($ctx, 'incr', 'default', array('_counter_for_demo_'));
    $r = $client->rd_call($ctx, 'get', 'default', array('_counter_for_demo_'));
    var_dump($r);

    for ($i = 0; $i < 500; $i++) {
        $lockKey = "foo";
        var_dump($client->lock($ctx, 'just a test', $lockKey));
        $client->unlock($ctx, 'just a test', $lockKey);
    }

    $t1 = microtime(TRUE);
    // game get unique name with len 3
    //for ($i = 0; $i < 2; $i ++) {
    //for ($i = 0; $i < 658; $i ++) {
    $allianceTags = array();
    for ($i = 0; $i < 50000000; $i ++) {
        //$name = $client->ping($ctx);
        $name = $client->gm_name3($ctx);
        echo "$i $name\n";
        if (isset($allianceTags[$name])) {
            throw new Exception("Dup name3: $name");
        }
        $allianceTags[$name] = TRUE;
        //usleep(10000);
        //sleep(1);
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

echo microtime(TRUE) - $t1, "\n";
