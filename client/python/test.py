#!/usr/bin/env python
#encoding=utf-8
'''
for quick debugging of fae
'''

import sys
import bson
import datetime
sys.path.append('../../servant/gen-py')
sys.path.append('/System/Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/site-packages')
from thrift.transport import TSocket
from thrift.protocol import TBinaryProtocol
from thrift.transport.TTransport import TTransportException
from thrift.Thrift import TApplicationException
from fun.rpc import FunServant
from fun.rpc.ttypes import TCacheMissed, TMemcacheData

t1 = datetime.datetime.now()
sock = TSocket.TSocket('localhost', 9001)
try:
    sock.open()
except TTransportException, e:
    print e
    sys.exit(1)
protocol = TBinaryProtocol.TBinaryProtocol(sock)

client = FunServant.Client(protocol)
ctx = FunServant.Context(reason='POST+/facebook/getPaymentRequestId/+34ca2cf6', rid='121')

def elapsed():
    global t1
    delta = datetime.datetime.now() - t1
    t1 = datetime.datetime.now()
    return '(' + str(delta.microseconds) + ' us)'

# ping
#=====
r = client.ping(ctx)
delta = datetime.datetime.now() - t1
print '[Client] ping received:', r, elapsed()

# mc
#=====
mcData = TMemcacheData()
mcData.data = 'world 世界'
print '[Client] mc_set received:', client.mc_set(ctx, 'default', 'hello', mcData, 120), elapsed()
print '[Client] mc_get received:', client.mc_get(ctx, 'default', 'hello'), elapsed()

try:
    print '[Client] mc_get hello-non-exist ->', client.mc_get(ctx, 'default', 'hello-non-exist'), elapsed()
except TCacheMissed:
    print 'cache missed'
except TApplicationException, e:
    print e
except Exception, e:
    print e, type(e)

# lc
#=====
print '[Client] lc_set received:', client.lc_set(ctx, 'lc_test_hello', 'abcdefg'), elapsed()
print '[Client] lc_get received:', client.lc_get(ctx, 'lc_test_hello'), elapsed()
print '[Client] lc_del received:', client.lc_del(ctx, 'lc_test_hello'), elapsed()

# mg
#=====
userDoc = {
    "name": "funky.gao",
    "gendar": "M",
    "abtype": {
        "payment": "a",
        "tutorial": "b",
    }
}
doc = bson.dumps(userDoc)
print 'doc:', doc
#print '[Client] mg_insert received:', client.mg_insert(ctx, 'db1', 'usertest', 0, doc), elapsed()

# my
#=====
print '[Client] my_query received:', client.my_query(ctx, 'UserShard', 'UserInfo', 1, 'SELECT * FROM UserInfo', ()), elapsed()
