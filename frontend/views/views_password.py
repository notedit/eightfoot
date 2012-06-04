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

from libshare import oocrpc
oocrpc.backend = settings.RPC


