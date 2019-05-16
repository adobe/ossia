#!/bin/bash

APP_NAME=ossia
DAEMON=/opt/$APP_NAME/bin/$APP_NAME
PID=$(pgrep -f $DAEMON)

if [[ ! -z "$PID" ]]; then
    if [[ -f /etc/redhat-release ]]; then
        if [[ "$1" = "0" ]]; then
            which systemctl &>/dev/null
            if [[ $? -eq 0 ]]; then
                systemctl stop $APP_NAME
            else
                service $APP_NAME stop
            fi
        fi
    elif [[ -f /etc/lsb-release ]]; then
        if [[ "$1" != "upgrade" ]]; then
            which systemctl &>/dev/null
            if [[ $? -eq 0 ]]; then
                systemctl stop $APP_NAME
            else
                stop $APP_NAME
            fi
        fi
    elif [[ -f /etc/os-release ]]; then
        source /etc/os-release
        if [[ $ID = "amzn" ]]; then
            if [[ "$1" = "0" ]]; then
                service $APP_NAME stop
            fi
        fi
    else
        kill -9 $PID
    fi
fi
