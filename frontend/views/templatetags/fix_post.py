# -*- coding: utf-8 -*-
# file: fix_post.py
# date: 20120716
# author: notedit



"""
一些辅助的tag
"""

import time
from django import template
from django.conf import settings

register = template.Library()

@register.filter
def format_datetime(post_time):
    return time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(post_time[0]))



