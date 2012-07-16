# -*- coding: utf-8 -*-
# date: 2012-05-29
# author: notedit

"""
Your time is limited, so don't waste it living someone else's life.
Don't be trapped by dogma - which is living with the results of other
people's thinking. Don't let the noise of other's opinions drown out
your own inner voice. And most important, have the courage to follow
your heart and intuition. They somehow already know what you truly
want to become. Everything else is secondary.

by Steve Jobs

"""

import os
from pprint import pprint
from django.conf import settings
from django.shortcuts import render_to_response
from django.http import HttpResponse

from libshare import oocrpc
from libshare import authutil
from libshare import strutil

RC = settings.RC
oocrpc.backend = settings.RPC

def index(req,page=1):
    """首页"""
    offset = (page-1)*25
    comm_dict = {}

    is_logined = authutil.is_logined(req)
    if is_logined:
    	curr_ukey = req.COOKIES.get('ukey')
    	follow_count = oocrpc.backend.GetFollowContentCount(curr_ukey)
    	follow_contents = oocrpc.backend.GetFollowContent({'Ukey':curr_ukey,'Offset':offset,'Limit':25})
    	pager = strutil.pager(page,follow_count,'/index/',per_page=25)
    	user_info = oocrpc.backend.GetUserInfo(curr_ukey)
    	comm_dict.update({'contents':follow_contents,'ukey':curr_ukey,'pager':pager,
    					  'is_logined':True,'user_info':user_info})
    else:
    	# hotest
    	hotest_count = oocrpc.backend.GetContentCount()  # to do 
    	hotest_contents = oocrpc.backend.GetHotestContent({'Offset':offset,'Limit':25}) # to do
    	pager = strutil.pager(page,hotest_count,'/index/',per_page=25)
    	comm_dict.update({'contents':hotest_contents,'pager':pager})

    pprint(comm_dict)

    return render_to_response('index.html',comm_dict)


def index_latest(req,page=1):
    offset = (page-1)*25
    comm_dict = {}
    newest_count = oocrpc.backend.GetContentCount()
    newest_contents = oocrpc.backend.GetLatestContent({'Offset':offset,'Limit':25})
    pager = strutil.paper(page,newest_acount,'/index/newest/',per_page=25)
    comm_dict.update({'newest_count':newest_count,'newest_contents':newest_contents,'pager':pager})

    is_logined = authutil.is_logined(req)
    if is_logined:
    	curr_ukey = req.COOKIES.get('ukey')
    	user_info = oocrpc.backend.GetUserInfo(curr_ukey)
    	comm_dict.update({'curr_ukey':curr_ukey,'user_info':user_info})
    return render_to_response('index_newest.html',comm_dict)


def index_hotest(req,page=1):
	offset = (page-1)*25
	comm_dict = {}
	hotest_count = oocrpc.backend.GetContentCount()
	hotest_contents = oocrpc.backend.GetHotestContent({'Offset':offset,'Limit':25})
	pager = strutil.paper(page,hotest_count,'/index/',per_page=25)
	comm_dict.update({'hotest_count':hotest_count,'hotest_contents':hotest_contents})

	is_logined = authutil.is_logined(req)
	if is_logined:
		curr_ukey = req.COOKIES.get('ukey')
		user_info = oocrpc.backend.GetUserInfo(curr_ukey)
		comm_dict.update({'curr_ukey':curr_ukey,'user_info':user_info})
	return render_to_response('index_hotest.html',comm_dict)


def test_rpc(req):
    username = "young man"
    username = oocrpc.backend.GetHelloWorld('hey young man') 
    return HttpResponse(username)


  ### Unittest  #################################################################

from django.utils import unittest
from django.test.client import Client 


class TestView(unittest.TestCase):

	def setUp(self):
		pass

	def test_index_hotest(self):
		pass

	def test_index_newest(self):
		pass

	def test_index(self):
		pass



