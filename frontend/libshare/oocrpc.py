# -*- coding: utf-8 -*-
# date: 2012-06-24
# author: notedit<notedit@gmail.com>

import socket
import traceback
import msgpack

"""
a simple connection client
"""

SOCKET_TIMEOUT = 10.0

default_read_buffer_size = 8192

class ConnectionError(Exception):
    pass

class RpcError(Exception):
    
    def __init__(self,message):
        self.message = message

    def __str__(self):
        return self.message

    def __repr__(self):
        return self.message

class BackendError(RpcError):
    pass


class Request(object):
    header = None
    body = None
    def __init__(self,method,args):
        self._operation = 1
        self._method = method
        self.body = args
        self.header = {'Operation':self._operation,
                        'Method':self._method}
    
    def encode_request(self):
        try:
            return msgpack.packb(self.header) + msgpack.packb(self.body)
        except Exception,ex:
            errstr = traceback.format_exc()
            raise RpcError('EncodeError:%s'%errstr)

class Response(object):
    
    def __init__(self,header,body):
        self.header = header
        self.reply = body

    @property
    def error(self):
        return self.header.get('Error')


class Connection(object):

    def __init__(self,host='localhost',port=9090,timeout=10.0):
        self.host = host
        self.port = port
        self.timeout = timeout
        self._conn = None

    def __del__(self):
        try:
            self.close()
        except:
            pass

    @property
    def conn(self):
        if self._conn:
            return self._conn
        try:
            sock = self.connect()
        except socket.timeout,e:
            raise ConnectionError('can not connect to %s:%d'%(self.host,self.port))
        except socket.error,e:
            raise ConnectionError('can not connect to %s:%d'%(self.host,self.port))
        self._conn = sock
        return self._conn

    def connect(self):
        sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        sock.settimeout(self.timeout)
        sock.connect((self.host,self.port))
        return sock

    def reconnect(self):
        self.close()
        try:
            sock = self.connect()
        except:
            raise ConnectionError('can not connect to %s:%d'%(self.host,self.port))
        self._conn = sock

    def close(self):
        if self._conn is None:
            return
        try:
            self._conn.close()
        except socket.error:
            pass
        self._conn = None

    def write_request(self,method,args):
        request = Request(method,args)
        data = request.encode_request()
        try:
            self.conn.sendall(data)
        except socket.error,ex:
            self.reconnect()
            self.conn.sendall(data)
        except socket.timeout,ex:
            self.reconnect()
            self.conn.sendall(data)

    def read_response(self):
        try:
            unpacker = msgpack.Unpacker()
            data = self._read_more()
            unpacker.feed(data)
            # read the header
            header = None
            while True:
                try:
                    header = unpacker.next()
                    break
                except StopIteration:
                    data = self._read_more()
                    unpacker.feed(data)
            # read the body
            body = None
            while True:
                try:
                    body = unpacker.next()
                    break
                except StopIteration:
                    data = self._read_more()
                    unpacker.feed(data)
            res = Response(header,body)
            return res
        except ConnectionError:
            raise
        except Exception:
            excstr = traceback.format_exc()
            raise RpcError(excstr)

    def _read_more(self):
        try:
            data = self.conn.recv(default_read_buffer_size)
            if not data:
                raise ConnectionError('unexpected EOF in read')
        except socket.error:
            raise ConnectionError('unexpected recv error')
        except socket.timeout:
            raise ConnectionError('unexpected timeout error')
        return data


class RpcClient(object):
    """rpc client"""
    def __init__(self,host='localhost',port=9090):
        self.host = host
        self.port = port
        self.conn = Connection(self.host,self.port)

    def __getattr__(self,funcname):
        func = lambda *args:self.__call__(funcname,*args)
        func.__name__ = funcname
        return func

    def __call__(self,method,*args):
        if len(args) > 1:
            raise RpcError('method should only have one parameter')
        elif len(args) == 0:
            arg = {}
        else:
            arg = args[0]
        self.conn.write_request(method,arg)
        res = self.conn.read_response()
        if res.error:
            raise RpcError(res.error)
        return res.reply

if __name__ == '__main__':
    client = RpcClient(host='localhost',port=9090)
    ret = client.Add({'A':7,'B':9})
    print 'Arith.Add',ret

    ret = client.Mul({'A':7,'B':8})
    print 'Arith.Mul',ret

    ret = client.Div({'A':9,'B':3})
    print 'Arith.Div',ret

    try:
        ret = client.NError({'A':9,'B':3})
    except RpcError,ex:
        print 'Arith.NError',str(ex)

    ret = client.SimpleValue(2)
    print 'SimpleValue',ret

    ret = client.SimpleValue(3)
    print 'SimpleValue',ret
