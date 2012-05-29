#!/usr/bin/env sh
# Author: notedit

PIDFILE="/tmp/eightfoot.pid"
APPPORT=8000

case $1 in
    start)
        gunicorn_django --workers=3 -k gevent -D  --pid $PIDFILE -b $APPPORT
        ;;
    stop)
        kill `cat $PIDFILE`
        ;;
    debug)
        gunicorn_django --workers=2 -k gevent -b $APPPORT
        ;;
    *)
        echo "Usage: ./ctlapp.sh start | stop | debug"
        ;;
esac
