# -*- coding: utf-8 -*-
# date: 2012-06-24
# author: notedit

import os
import traceback
from django.conf import settings
from django.shortcuts import render_to_response
from django.http import HttpResponse
from django.http import HttpResponseRedirect
from django.utils import simplejson

from libshare import oocrpc
from libshare import authutil

oocrpc.backend = settings.RPC


@authutil.ajax_api
def signup(req):
	if authutil.is_logined(req):
		return authutil.api_error(500,'已经是登陆状态')
	form = SignupForm(req.POST)
	if form.is_valid():
		email = form.cleaned_data['email'].encode('utf-8')
		password = form.cleaned_data['password'].encode('utf-8')
		nickname = form.cleaned_data['nickname'].encode('utf-8')
		try:
			ukey = oocrpc.backend.RegisterUser({'Nickname':nickname,'Email':email,'Password':password})
		except oocrpc.RpcError, ex:
			if ex.message.startswith('EmailError'):
				return HttpResponse(simplejson.dumps({'emailerror':'邮箱已经存在'}))
			elif ex.message.startswith('NicknameError'):
				return HttpResponse(simplejson.dumps({'nicknameerror':'用户昵称重复'}))
			else:
				return HttpResponse(simplejson.dumps({'internalerror':'服务器开小差了'}))
		resp = HttpResponse(simplejson.dumps({'ok':1}))
		timeout = 3600*24*30*6
		authutil.set_logined(req,resp,ukey,timeout)
		return resp
	else:
		return HttpResponse(simplejson.dumps({'requesterror':'参数格式错误'}))


@authutil.ajax_api
def login(req):
	if authutil.is_logined(req):
		return authutil.api_error(500,'已经是登陆状态')
	form = LoginForm(req.POST)
	if form.is_valid():
		email = form.cleaned_data['email'].encode('utf-8')
		password = form.cleaned_data['password'].encode('utf-8')
		remember = form.cleaned_data['remember']
		try:
			ukey = oocrpc.backend.Login({'Email':email,'Password':password})
			userinfo = oocrpc.backend.GetUserInfo(ukey)
		except oocrpc.RpcError,ex:
			if ex.message.startswith('EmailError'):
				return HttpResponse(simplejson.dumps({'emailerror':'邮箱不存在'}))
			elif ex.message.startswith('PasswordError'):
				return HttpResponse(simplejson.dumps({'passworderror':'密码不正确'}))
			else:
				# todo add log
				return HttpResponse(simplejson.dumps({'internalerror':'服务器开小差了'}))
		retdict = {
			'ukey':ukey,
			'userinfo':userinfo
		}
		resp = HttpResponse(retdict)
		if remember:
			timeout = 3600*24*30*6  # half of a year
		else:
			timeout = None
		authutil.set_logined(req,resp,ukey,timeout)
		return resp
	else:
		return HttpResponse(simplejson.dumps({'requesterror':'参数格式错误'}))


@authutil.ajax_api
def logout(req):
	resp = HttpResponse('OK')
	authutil.set_logout(req,resp)
	return resp

@authutil.ajax_api
def verify_nickname(req):
	form = VerifyNicknameForm(req.POST)
	if form.is_valid():
		nickname = form.cleaned_data['nickname']
		try:
			ukey = oocrpc.backend.VerifyNickname(nickname.encode('utf-8'))
		except oocrpc.RpcError,ex:
			return authutil.api_error(500,u'服务器开小差了')
		if ukey:
			return HttpResponse(simplejson.dumps({'is_valid':False}))
		else:
			return HttpResponse(simplejson.dumps({'is_valid':True}))
	else:
		return authutil.api_error(500,u'参数格式错误:%s'%form.errors.as_text())

@authutil.ajax_api
def verify_email(req):
	form = VerifyEmailForm(req.POST)
	if form.is_valid():
		email = form.cleaned_data['email']
		try:
			"backend-service"
			ukey = oocrpc.backend.VerifyEmail(email.encode('utf-8'))
		except oocrpc.RpcError,ex:
			return authutil.api_error(500,u'服务器开小差了')
		if ukey:
			return HttpResponse(simplejson.dumps({'is_valid':False}))
		else:
			return HttpResponse(simplejson.dumps({'is_valid':True}))
	else:
		return authutil.api_error(500,u'参数格式错误:%s'%form.errors.as_text())




### Unittest  #################################################################

from django.utils import unittest
from django.test.client import Client 


class TestView(unittest.TestCase):

	def setUp(self):
		pass

	def test_signup(self):
		pass

	def test_login(self):
		pass

	def test_logout(self):
		pass

	def test_verify_nickname(self):
		pass

	def test_verify_email(self):
		pass
