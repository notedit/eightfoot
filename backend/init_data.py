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

    for idx in xrange(0,99):
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
    t = insertModel('tag',{
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

    for idx in xrange(101,200):
        t.append({
            'name':'tagname%03d' % idx,
            'introduction':'this is %03d' % idx,
            'date_create':date_str,
            'author_ukey':'user%02d' % (idx%9 + 1),
            'url_code':''
            })
    print str(t)

def all_content(date_str):
    dt=time.mktime(time.strptime(date_str,'%Y-%m-%d %H:%M:%S'))
    c = insertModel('content',{
        'title':strType,
        'author_ukey':strType,
        'last_modify_ukey':strType,
        'last_reply_ukey':strType,
        'body':strType,
        'url':strType,
        'atype':strType,
        'date_create':strType
        })
    for idx in xrange(1,100):
        c.append({
            'title':'我们纷纷表示压力很大，title %02d' % idx,
            'author_ukey':'user%02d' % idx,
            'last_modify_ukey':'user%02d' % idx,
            'last_reply_ukey':'user%02d' % idx,
            'body':'this is a  beautiful body  body %02d' % idx,
            'url':'http://www.baidu.com',
            'atype':'content',
            'date_create':time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(dt+idx*2))
            })

    for idx in xrange(100,200):
        c.append({
            'title':'我们纷纷表示压力很大，title %03d' % idx,
            'author_ukey':'user%02d' % (idx-100) ,
            'last_modify_ukey':'user%02d' % (idx-100),
            'last_reply_ukey':'user%02d' % (idx-100),
            'body':'this is a  beautiful body  body %03d' % idx,
            'url':'http://www.baidu.com',
            'atype':'content',
            'date_create':time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(dt+idx*2))
            })
    print str(c)

def all_tag_map():
    tm = insertModel('tag_map',{
        'tag_id':intType,
        'content_id':intType
        })
    for idx in xrange(1,5):
        tids = [1,2,3]
        for tid in tids:
            tm.append({
                'tag_id':tid,
                'content_id':idx
                })
    tagids = [i for i in xrange(1,20)]
    for idx in xrange(5,200):
        tids = random.sample(tagids,3)
        for tid in tids:
            tm.append({
                'tag_id':tid,
                'content_id':idx
            })
    print str(tm)

if __name__ == '__main__':
    all_user()
    all_tag('2012-05-13 17:35:45')
    all_content('2012-03-13 17:35:45')
    all_tag_map()
