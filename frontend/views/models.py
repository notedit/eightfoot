# -*- coding: utf-8 -*-
# date: 2012-05-18
# author: notedit

"""
记录一些公用的东西，比如数据操作的封装，
比如数据库的读操作
"""

from django.conf import settings
from django.template import Context,loader
from django.http import HttpResponse
from django.utils.http import cookie_date





class LoginedCookieMiddleware(object):

    def process_response(self,request,response):
        if not hasattr(request,'session'):
            return response
        #if request.session.has_key('is_logined') and request.session['is_logined']:
        if request.session.get('is_logined'):
            max_age = request.session.get_expiry_age()
            expires_time = time.time() + max_age
            expires = cookie_date(expires_time)
            response.set_cookie('is_logined',
                            'TRUE',
                            max_age=max_age,
                            expires=expires, domain=settings.SESSION_COOKIE_DOMAIN,
                            path=settings.SESSION_COOKIE_PATH,
                            secure=settings.SESSION_COOKIE_SECURE or None)
            #response.set_cookie('is_logined','TRUE',max_age=request.session.get_expiry_age())
        else:
            response.delete_cookie('is_logined')
        return response




