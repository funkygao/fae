#!/usr/bin/env python
#encoding=utf-8
'''
stress test for dryrun of fae: ping
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
PINGS_PER_CLIENT = 2

def ping(n):
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
        client.ping(ctx)

def main():    
    t1 = datetime.datetime.now()
    pool = multiprocessing.Pool(processes=CONCURRENCY)
    for i in xrange(SESSIONS):
        pool.apply(ping, (PINGS_PER_CLIENT, ))
    pool.close()
    pool.join()

    print PINGS_PER_CLIENT*SESSIONS, 'called'
    print datetime.datetime.now() - t1

if __name__ == '__main__':
    main()
