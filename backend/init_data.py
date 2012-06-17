#!/usr/bin/env python2.7
# -*- coding: utf-8 -*-
# date: 20120526

import os
import sys
import time
import socket
import random
import itertools
import hashlib



HOSTNAME=socket.gethostname()
CURPATH=os.path.normpath(os.path.join(os.getcwd(),os.path.dirname(__file__)))
int2fmt = lambda i: time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(i))

salt=lambda length:''.join([chr(random.choice(range(97,123))) for x in range(length)])

SQL_SEQ_SETVAL="""SELECT setval(%(seqkey)s,%(seqval)d);"""

class typeBase:
    vtype = 'Base'
    def __init__(self, val):
        self.val = val

    def __str__(self):
        raise Exception("NotImplementError")


class strType(typeBase):
    vtype = "strType"
    def __str__(self):
        return "'%s'" % self.val

class intType(typeBase):
    vtype = "intType"
    def __str__(self):
        return str(self.val)

class insertModel:
    def __init__(self, table, fields):
        self.fields={}
        self.values=[]
        self.table=''
        assert str == type(table), "table name should be str type"
        assert dict == type(fields), "fields should be dict type"
        self.fields = fields
        self.keys = fields.keys()
        self.table = table

    def append(self, record):
        lostfields=filter(lambda k:k not in record.keys(),self.keys)
        assert lostfields==[],"Lost field=%s"%repr(lostfields)
        tmpvals = """( %s )""" % ', '.join(map(lambda k: str(self.fields[k](record[k])), self.keys))
        self.values.append(tmpvals)

    def __str__(self):
        return """INSERT INTO %s (%s) VALUES %s;""" % ( self.table, ', '.join(self.keys), ',\n'.join(self.values))



def all_user():
    u = insertModel('user_info',{
        'ukey':strType,
        'nickname':strType,
        'pic_small':strType,
        'pic_big':strType,
        'status':intType,
        'introduction':strType,
        })

    for idx in range(99):
        u.append({
            'ukey': 'user%02d' % (idx+1),
            'nickname':'nick%02d' % (idx+1),
            'pic_small':'pic_small',
            'pic_big':'pic_big',
            'status':1,
            'introduction':'introduction.....'
            })
    print str(u)

    p = insertModel('passwd',{
        'ukey':strType,
        'email':strType,
        'salt':strType,
        'password':strType
        })

    for idx in range(99):
        saltstr = salt(10)
        passstr = 'pass%02d' % (idx+1)
        p.append({
            'ukey':'user%02d' % (idx+1),
            'email':'user%02d@gmail.com' % (idx+1),
            'salt':saltstr,
            'password':hashlib.sha1(saltstr+'pass%02d' % (idx+1)).hexdigest(),
            })
    print str(p)

def all_tag(date_str):
    t = insertModel({
        'name':strType,
        'introduction':strType,
        'date_create':strType,
        'author_ukey':strType,
        'url_code':strType
        })

    for idx in xrange(1,100):
        t.append({
            'name':'tagname%02d' % idx,
            'introduction':'this is %02d' % idx,
            'date_create':date_str,
            'author_ukey':'user%02d' % idx,
            'url_code':''
            })
    print str(t)

    for idx in xrange(101,200):
        t.append({
            
            })

def all_content(date_str):
    c = insertModel('content',{
        'title':strType,
        'author_ukey':strType,
        'last_modify_ukey':strType,
        'last_reply_ukey':strType,
        'body':strType,
        })
    for idx in range(99):
        c.append({
            'title':'我们纷纷表示压力很大，title %02d' % (idx+1),
            'author_ukey':'user%02d' % (idx+1),
            'last_modify_ukey':'user%02d' % (idx+1),
            'last_reply_ukey':'user%02d' % (idx+1),
            'body':'this is a  beautiful body  body %02d' % (idx+1)
            })
    print str(c)

    for idx in xrange(100,199):
        c.append({
            'title':'我们纷纷表示压力很大，title %03d' % idx,
            'author_ukey':'user%03d' % idx,
            'last_modify_ukey':'user%03d' % idx,
            'last_reply_ukey':'user%03d' % idx,
            'body':'this is a  beautiful body  body %03d' % idx)
            })
    print str(c)



if __name__ == '__main__':
    all_user()
