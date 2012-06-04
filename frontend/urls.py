# -*- coding: utf-8 -*-
# date: 2012-05-18
# author: notedit

from django.conf.urls import patterns, include, url
from django.views.generic.simple import direct_to_template
from django.conf import settings

urlpatterns = patterns('',
    url(r'^api/',       include('api.urls')),
)

if settings.HOSTNAME in ('notedit'):
    urlpatterns += patterns('',
        url(r'^css/(?P<path>.*)$','django.views.static.serve',
            {'document_root':os.path.join(settings.CURR_PATH,'static','css')}),
         url(r'^js/(?P<path>.*)$','django.views.static.serve',
            {'document_root':os.path.join(settings.CURR_PATH,'static','js')}),

    )


urlpatterns += patterns('views.views_index',
        url(r'^$','index'),
        url(r'^test_rpc/$','test_rpc'),
        )


