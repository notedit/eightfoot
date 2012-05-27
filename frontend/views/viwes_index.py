# -*- coding: utf-8 -*-
# date: 2012-05-18
# author: notedit

import os

from django.conf import settings
from django.shortcuts import render_to_response
from django.http import HttpResponse

oocrpc = settings.RPCCLIENT

def index(req):
    return HttpResponse("hello world, this is the index")

def test_rpc(req):
    username = oocrpc.getUserName('hey young man')
    return HttpResponse('hey ' + username)



