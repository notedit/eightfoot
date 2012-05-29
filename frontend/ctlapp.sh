#!/usr/bin/env sh
# Author: notedit

PIDFILE="/tmp/eightfoot.pid"
APPADDR="127.0.0.1:8000"

case $1 in
    start)
        gunicorn_django --workers=3 -k gevent -D  --pid $PIDFILE
        ;;
    stop)
        kill -INT `cat $PIDFILE`
        ;;
    debug)
        gunicorn_django --workers=2  --pid $PIDFILE
        ;;
    *)
        echo "Usage: ./ctlapp.sh start | stop | debug"
        ;;
esac
