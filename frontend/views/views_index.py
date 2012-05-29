# -*- coding: utf-8 -*-
# date: 2012-05-29
# author: notedit

import os

from django.conf import settings
from django.shortcuts import render_to_response
from django.http import HttpResponse

from libshare import oocrpc

def index(req):
    return HttpResponse("hello world, this is the index")

def test_rpc(req):
    username = "young man"
    username = oocrpc.backend.getHelloWorld('hey young man')
    return HttpResponse(username)



