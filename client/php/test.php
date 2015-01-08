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
    $sock->setDebug(0);
    $sock->setSendTimeout(4000);
    $sock->setRecvTimeout(4000);
    $sock->setNumRetries(1);
    $transport = new TBufferedTransport($sock, 1024, 1024);
    $protocol = new TBinaryProtocol($transport);

    // get our client
    $client = new FunServantClient($protocol);
    $transport->open();

    $ctx = new Context(array('reason' => "phptest", 'rid' => '125'));

    // ping
    $return = $client->ping($ctx);
    echo "[Client] ping received: ", $return, "\n";

    // id.next
    echo "[Client] id_next received:", $client->id_next($ctx), "\n";
    echo "[Client] id_next received:", $client->id_next($ctx), "\n";
    echo "[Client] id_next_with_tag received:", $client->id_next_with_tag($ctx, 5), "\n";
    echo "[Client] id_next_with_tag received:", $client->id_next_with_tag($ctx, 10), "\n";
    $id = $client->id_next_with_tag($ctx, 18);
    list($ts, $tag, $wid, $seq) = $client->id_decode($ctx, $id);
    echo "$ts $tag $wid $seq\n";

    // lc
    echo '[Client] lc_set received: ', $client->lc_set($ctx, 'hello-php-lc', 'world 世界'), "\n";
    echo '[Client] lc_get received: ', $client->lc_get($ctx, 'hello-php-lc'), "\n";
    echo '[Client] lc_del received: ', $client->lc_del($ctx, 'hello-php-lc'), "\n";

    // my.query
    if (1) {
        for ($i=0; $i<5; $i++) {
            $rows = $client->my_query($ctx, 'UserShard', 'UserInfo', 1, 'SELECT * FROM UserInfo', array(), '');
            echo $rows->rowsAffected, ':rowsAffected, ', $rows->lastInsertId, ':lastInsertId, rows:', PHP_EOL;
            print_r($rows);
        }
    }

    // mc
    if (0) {
        $mcData = new TMemcacheData();
        $mcData->data = 'world 世界';
        echo '[Client] mc_set received: ', $client->mc_set($ctx, 'default', 'hello-php', $mcData, 120), "\n";
        echo '[Client] mc_get received: ', print_r($client->mc_get($ctx, 'default', 'hello-php')), "\n";
        $mcData->data = 0;
        echo '[Client] mc_add received: ', $client->mc_add($ctx, 'default', 'test:counter:uid', $mcData, 3500), "\n";
        echo '[Client] mc_inc received: ', $client->mc_increment($ctx, 'default', 'test:counter:uid', 7), "\n";
        try {
            echo '[Client] mc_get hello-non-exist received: ', $client->mc_get($ctx, 'default', 'hello-non-exist'), "\n";
        } catch (TCacheMissed $ex) {
            echo 'mc error: ', $ex->getMessage(), "\n";
        }
    }

    // mg.insert
    if (0) {
        $doc = array(
            "name" => "funky.php",
            "gendar" => "M",
            "abtype" => array(
                "payment" => "a",
                "tutorial" => "b",
            )
        );
        echo '[Client] mg_insert received: ', $client->mg_insert($ctx, 'db1', 'usertest', 0, 
            bson_encode($doc)), "\n";

        // mg.inserts
        $docs = array();
        $docs[] = bson_encode($doc);
        $docs[] = bson_encode($doc);
        echo '[Client] mg_inserts received: ', $client->mg_inserts($ctx, 'db1', 'usertest2', 0, 
            $docs), "\n";

        // mg.findOne
        try {
            $idmap = $client->mg_find_one($ctx, 'default', 'idmap', 0,
                bson_encode(array('snsid' => '100003391571259')), bson_encode(''));
            echo "[Client] mg_find_one received: \n";
            print_r(bson_decode($idmap));
        } catch (TMongoMissed $ex) {
            echo $ex->getMessage(), "\n";
        }

        // mg.count
        echo "[Client] mg_count received:", $client->mg_count($ctx, 'default', 'idmap', 0,
            bson_encode(array('uid' => array('$gte' => 1)))), "\n";
        echo "[Client] mg_count received:", $client->mg_count($ctx, 'default', 'idmap', 0,
            bson_encode(array('uid' => array('$gte' => 100000)))), "\n";

        // mg.findAll
        echo "[Client] mg_find_all received: \n";
        try {
            $docs = $client->mg_find_all($ctx, 'default', 'idmap', 0,
                bson_encode(array('uid' => array('$gte' => 1))), bson_encode(array()),
                0, 0, array());
            $r = array();
            foreach ($docs as $doc) {
                $r[] = bson_decode($doc);
            }
            print_r($r);
        } catch (TProtocolException $ex) {
            print_r($ex);
        }

        // mg.findAndModify
        $r = $client->mg_find_and_modify($ctx, 'default', 'idsquence', 0,
            bson_encode(array('table_name' => 'idMap')),
            bson_encode(array('$inc' => array('value' => 1))),
            true,
            false,
            true);
        $val = bson_decode($r);
        echo "[Client] mg.find_and_modify received: ", $val['value'], "\n";
    }

    $transport->close();
} catch (Exception $tx) {
    print 'Something went wrong: ' . $tx->getMessage() . "\n";
}

