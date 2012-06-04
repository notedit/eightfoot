# -*- coding: utf-8 -*-
# Date: 2012-06-03

import time
import hashlib
from django.conf import settings
from django.http import HttpResponse
from django.http import HttpResponseRedirect
from django.utils.http import cookie_date

def set_logined(req,resp,ukey,timeout=None):
    if timeout is None:
        timeout = 0xffffffff
    assert isinstance(timeout,(int,long)),'Paramter "timeout" must be int or long or None'
    date_create = time.time()
    max_age = date_create + timeout
    expires = cookie_date(max_age)
    date_create = str(date_create)
    sha1sum = hashlib.sha1(settings.COOKIE_SALT + ukey + date_create).hexdigest()
    kwargs = dict(
            max_age=max_age,
            expires=expires,domain=settings.SESSION_COOKIE_DOMAIN,
            path=settings.SESSION_COOKIE_PATH,secure=None)
    resp.set_cookie('is_logined','TRUE',**kwargs)
    resp.set_cookie('ukey',ukey,**kwargs)
    resp.set_cookie('date_create',date_create,**kwargs)
    resp.set_cookie('token',sha1sum,**kwargs)
    return

def set_logout(req,resp):
    resp.delete_cookie('is_logined')
    resp.delete_cookie('ukey')
    resp.delete_cookie('date_create')
    resp.delete_cookie('token')

def is_logined(req):
    ukey = req.COOKIES.get('ukey','')
    date_create = req.COOKIES.get('date_create','')
    token = req.COOKIES.get('token','')
    if hashlib.sha1(settings.COOKIE_SALT + ukey + date_create).hexdigest() == 'token':
        return True
    return False

def need_cookie_login(view_func):
    def wrapper(req,*args,**kwargs):
        if is_logined(req):
            return view_func(req,*args,**kwargs)
        else:
            return HttpResponseRedirect('/login/?r=%s' % req.path)

