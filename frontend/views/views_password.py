# -*- coding: utf-8 -*-
# date: 2012-05-29
# author: notedit

import os
import re
from urllib import urlencode,urlopen
from contextlib import closing
from django.conf import settings
from django.shortcuts import render_to_response
from django.http import HttpResponse
from django.http import HttpResponseRedirect

from pygithub3 import Github
from libshare import oocrpc
oocrpc.backend = settings.RPC

OAUTH_TOKEN_RE = re.compile(r'''<OAuth><token_type>(?p<token_type>.*?)</token_type><access_token>(?P<access_token>.*?)</access_token></OAuth>''')

def login_oauth(req):
    data = [('client_id',''),('redirect_url',''),('scope','user,repo,gist')]
    oauth_url = 'https://github.com/login/oauth/authorie?%s' % urlencode(data)
    return HttpResponseRedirect(oauth_url)

def oauth_callback(req,code=None):
    if not code:
        return HttpResponseRedirect('/error') # todo 
    oauth_token_url = 'https://github.com/login/oauth/access_token'
    data = [('client_id',''),('redirect_url',''),
            ('client_secret',''),('code',code)]
    access_token = None
    with closing(urlopen(oauth_token_url,data=urlencode(data).encode('utf-8'))) as f:
        ret = f.read() 
        gdict = dict(OAUTH_TOKEN_RE.findall(ret))
        token_type = gdict.get('token_type')
        access_token = gdict.get('access_token')
        if not access_token:
            return HttpResponseRedirect('/error')

    gh = Github(token=access_token)
    user = gh.users.get()
    
