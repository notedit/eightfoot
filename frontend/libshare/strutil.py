#!/usr/bin/evn python
# -*- coding: utf-8 -*-
# author: notedit
# date: 20120226

import hashlib
import random
from urlparse import urlparse,urlunparse

userid = lambda x : ''.join([random.choice('0123456789') for i in xrange(x)])
requestid = userid

sha1sum = lambda data: hashlib.sha1(data).hexdigest()
md5sum = lambda data: hashlib.md5(data).hexdigest()


def sanitize_url(url):
    if not url:
        return None
    url = url.strip().lower()
    try:
        u = urlparse(url)
        if not u.scheme:
            url = 'http://' + url
            u = urlparse(url)
    except:
        return None
    return url


def get_page_num(req,**kwargs):
    retdict={}
    for (k,v) in kwargs.items():
        try:
            new_value=int(req.GET.get(k,v))
            if new_value<=0:
                retdict[k]=v
            else:
                retdict[k]=new_value
        except TypeError,ex:
            retdict[k]=v
        except ValueError,ex:
            retdict[k]=v
    return retdict


def pager_bootstrap(page_num,tol_count,base_url,per_page=25,nav_len=10):
    '''
    分页相关的代码
    page_num : 当前页码
    tol_count: 对象总数
    base_url: 当前页面的url
    per_page: 每页显示对象数
    nav_len: 显示的分页个数
    '''
    qd = {'base_url':base_url,'num':''}
    if tol_count <= per_page:
        return ''
    if page_num == 1:
        prev_li_str = '<li class="disable"><a href="#">&larr;</a></li>\n'
    else:
        qd['num'] = page_num - 1
        prev_li_str = '<li><a href="%(base_url)s/page/%(num)s/">&larr;</a></li>\n' % qd

    page_count = (tol_count+per_page-1)/per_page
    if page_num == page_count:
        next_li_str = '<li class="disable"><a href="#">&rarr;</a></li>\n'
    else:
        qd['num'] = page_num + 1
        next_li_str = '<li><a href="%(base_url)s/page/%(num)s/">&rarr;</a></li>\n' % qd

    nav_len = nav_len / 2
    page_start = page_num - nav_len if (page_num - nav_len) > 1 else 1
    page_end = page_num + nav_len + 1 if page_num + nav_len <= page_count else page_count + 1
    page_range = range(page_start,page_end)
    if len(page_range) > 0 and page_range[0] != 1:
        page_range.insert(0,-1)
    if len(page_range) > 0 and page_range[-1] != page_count:
        page_range.append(-1)

    middle_line = []
    for num in page_range:
        if num == -1:
            middle_line.append('<li class="disable"><a href="#">...</a></li>\n')
        elif num == page_num:
            middle_line.append('<li class="active"><a href="#">%d</a></li>\n'%num)
        else:
            qd['num'] = num
            qd['shownum'] = num
            middle_line.append('<li><a href="%(base_url)s/page/%(num)s/">%(shownum)s</a></li>\n'%qd)

    html = '<div class="pagination-centered pagination">\n<ul>%s%s%s</ul>\n</div>\n'%(
                prev_li_str,
                ''.join(middle_line),
                next_li_str,
            )
    return html


# for rao
def pager(page_num,tol_count,base_url,per_page=25,nav_len=10):
    '''
    分页相关的代码
    page_num : 当前页码
    tol_count: 对象总数
    base_url: 当前页面的url
    per_page: 每页显示对象数
    nav_len: 显示的分页个数
    '''
    qd = {'base_url':base_url,'num':''}
    if tol_count <= per_page:
        return ''
    if page_num == 1:
        prev_li_str = '<li class="unavailable"><a href="javascript:void(0)">&laquo;</a></li>\n'
    else:
        qd['num'] = page_num - 1
        prev_li_str = '<li><a href="%(base_url)s/page/%(num)s/">&larr;</a></li>\n' % qd

    page_count = (tol_count+per_page-1)/per_page
    if page_num == page_count:
        next_li_str = '<li class="unavailable"><a href="javascript:void(0)">&laquo;</a></li>\n'
    else:
        qd['num'] = page_num + 1
        next_li_str = '<li><a href="%(base_url)s/page/%(num)s/">&raquo;</a></li>\n' % qd

    nav_len = nav_len / 2
    page_start = page_num - nav_len if (page_num - nav_len) > 1 else 1
    page_end = page_num + nav_len + 1 if page_num + nav_len <= page_count else page_count + 1
    page_range = range(page_start,page_end)
    if len(page_range) > 0 and page_range[0] != 1:
        page_range.insert(0,-1)
    if len(page_range) > 0 and page_range[-1] != page_count:
        page_range.append(-1)

    middle_line = []
    for num in page_range:
        if num == -1:
            middle_line.append('<li class="unavailable"><a href="#">...</a></li>\n')
        elif num == page_num:
            qd.update({'num':num,'shownum':num})
            middle_line.append('<li class="current"><a href="%(base_url)s/page/%(num)s/">%(shownum)s</a></li>\n'%qd)
        else:
            qd['num'] = num
            qd['shownum'] = num
            middle_line.append('<li><a href="%(base_url)s/page/%(num)s/">%(shownum)s</a></li>\n'%qd)

    html = '<div class="pagination-centered pagination">\n<ul>%s%s%s</ul>\n</div>\n'%(
                prev_li_str,
                ''.join(middle_line),
                next_li_str,
            )
    return html 
