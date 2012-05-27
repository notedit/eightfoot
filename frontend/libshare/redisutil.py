#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: notedit
# date: 20120301

import redis
from redis import ConnectionPool
from singleton import Singleton

class Redis(Singleton):

    db = None

    @staticmethod
    def create(**kwargs):
        r = redis.Redis(connection_pool=ConnectionPool(**kwargs))
        Redis.db = r
