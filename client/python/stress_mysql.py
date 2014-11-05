#!/usr/bin/env python
#encoding=utf-8
'''
stress test for mysql query with pk hit
'''

import sys
import datetime
import multiprocessing
sys.path.append('../../servant/gen-py')
sys.path.append('/System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/site-packages')
from thrift.transport import TSocket
from thrift.protocol import TBinaryProtocol
from thrift.transport.TTransport import TTransportException
from fun.rpc import FunServant

# config
CONCURRENCY = 100
SESSIONS = 1000 * 1000
QUERY_PER_CLIENT = 20

def mysql_query(n):
    #sock = TSocket.TSocket('127.0.0.1', 9001)
    sock = TSocket.TSocket('192.168.22.160', 9001)
    try:
        sock.open()
    except TTransportException, e:
        print e
        sys.exit(1)

    protocol = TBinaryProtocol.TBinaryProtocol(sock)
    client = FunServant.Client(protocol)
    ctx = FunServant.Context(caller='POST+/facebook/getPaymentRequestId/+34ca2cf6')
    for i in xrange(n):
        client.my_query(ctx, 'UserShard', 'UserInfo', 1, 'select * from UserInfo where uid=?', [1])

def main():    
    t1 = datetime.datetime.now()
    pool = multiprocessing.Pool(processes=CONCURRENCY)
    for i in xrange(SESSIONS):
        pool.apply(mysql_query, (QUERY_PER_CLIENT, ))
    pool.close()
    pool.join()

    print QUERY_PER_CLIENT*SESSIONS, 'called'
    print datetime.datetime.now() - t1

if __name__ == '__main__':
    main()
