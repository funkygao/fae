<?php
/**
 * Bootstrap of fae(Fun App Engine).
 *
 */

$GLOBALS['THRIFT_ROOT'] = __DIR__ . '/thrift';
$GLOBALS['FAE_ROOT'] = __DIR__ . '/fae'; // use fae 'build.sh -php' to generate php stubs
require_once $GLOBALS['THRIFT_ROOT'] . '/Thrift.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Base/TBase.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Exception/TException.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Exception/TProtocolException.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Exception/TApplicationException.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Exception/TTransportException.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Protocol/TBinaryProtocol.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Protocol/TBinaryProtocolAccelerated.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/StringFunc/TStringFunc.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/StringFunc/Core.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Factory/TStringFuncFactory.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Type/TType.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Type/TMessageType.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Transport/TSocket.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Transport/TSocketPool.php';
require_once $GLOBALS['THRIFT_ROOT'] . '/Transport/TBufferedTransport.php';
require_once $GLOBALS['FAE_ROOT'] . '/FunServant.php';
require_once $GLOBALS['FAE_ROOT'] . '/Types.php';

/**
 * A helper class to talk to fae.
 */
final class Fae {

    /**
     * @var fun\rpc\FunServantClient
     */
    private static $_client;

    /**
     * @var Thrift\Transport\TBufferedTransport
     */
    private static $_transport;

    /**
     * @var func\rpc\Context Context Stay the same within a php request lifespan
     */
    private static $_ctx = NULL;

    private static $_connected = FALSE;

    public function __destruct() {
        self::$_transport->close();
    }
    
    /**
     * Get a fae client connection.     
     *
     * @throws InvalidArgumentException
     * @return fun\rpc\FunServantClient
     */
    public static function client(array $hosts, array $ports, array $config) {
        static $instance = NULL;
        if (NULL === $instance) {
            $instance = new self(); // only to trigger __destruct

            $config += array(
                'send_timeout' => 4000, // ms
                'recv_timeout' => 4000, // ms
                'write_buffer' => 2048, // byte
                'read_buffer' => 2048,  // byte
                'retries' => 1,
            );

            // TSocketPool will auto try servers in turn if previous fae server is down
            $sock = new Thrift\Transport\TSocketPool($hosts, $ports);
            $sock->setDebug(0);
            $sock->setRandomize(TRUE); // simple fae load balance, fae is stateless
            $sock->setSendTimeout($config['send_timeout']); // TODO send timeout is also connection timeout
            $sock->setRecvTimeout($config['recv_timeout']);
            $sock->setNumRetries($config['retries']);
            self::$_transport = new Thrift\Transport\TBufferedTransport($sock, $config['read_buffer'],
                $config['write_buffer']);
            $protocol = new Thrift\Protocol\TBinaryProtocol(self::$_transport);
            self::$_client = new fun\rpc\FunServantClient($protocol);
            self::$_transport->open(); // eager connection
            self::$_connected = TRUE;
        }

        return self::$_client;
    }    

    /**
     * Create a fae context.
     *
     * @return fun\rpc\Context
     */
    public static function ctx($reason, $uid) {
        if (self::$_ctx == NULL) {
            $ctxInfo = array(
                // required
                'rid' => hexdec(uniqid());, // PHP request id
                'reason' => $reason,
                'uid' => $uid,
            );

            self::$_ctx = new fun\rpc\Context($ctxInfo);
        }
        
        return self::$_ctx;
    }

}
