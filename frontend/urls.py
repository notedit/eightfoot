# -*- coding: utf-8 -*-
# date: 2012-05-18
# author: notedit

import os
from django.conf.urls import patterns, include, url
from django.views.generic.simple import direct_to_template
from django.conf import settings

urlpatterns = patterns('',
    url(r'^api/',       include('api.urls')),
)

if settings.HOSTNAME in ('notedit','localhost'):
    urlpatterns += patterns('',
        url(r'^css/(?P<path>.*)$','django.views.static.serve',
            {'document_root':os.path.join(settings.CURR_PATH,'static','css')}),
         url(r'^js/(?P<path>.*)$','django.views.static.serve',
            {'document_root':os.path.join(settings.CURR_PATH,'static','js')}),
         url(r'^images/(?P<path>.*)$','django.views.static.serve',
             {'document_root':os.path.join(settings.CURR_PATH,'static','image')}),

    )


urlpatterns += patterns('views.views_index',
        url(r'^$','index'),
        url(r'^/index/(?P<page>\d{1,10})/$','index'),
        url(r'^/index/latest/(?P<page>\d{1,10})?/?','index_latest'),
        url(r'/index/hotest/(?P<page>\d{1,10})?/?','index_hotest'),
        url(r'^test_rpc/$','test_rpc'),
        )


