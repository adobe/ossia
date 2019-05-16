#!/bin/bash

APP_NAME=ossia
BASE_DIR=/opt/$APP_NAME
LOG_DIR=$BASE_DIR/log
DATA_DIR=$BASE_DIR/db

function remove_traces {
    test -d $BASE_DIR && rm -rf $BASE_DIR
    id $APP_NAME &>/dev/null
    if [[ $? -ne 1 ]]; then
          userdel $APP_NAME
    fi
}

function disable_systemd {
    systemctl disable $APP_NAME
    rm -f /lib/systemd/system/$APP_NAME.service
}

function disable_update_rcd {
    update-rc.d -f $APP_NAME remove
    rm -f /etc/init/$APP_NAME.conf
}

function disable_chkconfig {
    chkconfig --del $APP_NAME
    rm -f /etc/init.d/$APP_NAME
}

if [[ -f /etc/redhat-release ]]; then
    if [[ "$1" = "0" ]]; then
	     rm -f /etc/default/$APP_NAME
       which systemctl &>/dev/null
	     if [[ $? -eq 0 ]]; then
	        disable_systemd
	     else
	        disable_chkconfig
	     fi
    fi
elif [[ -f /etc/lsb-release ]]; then
    if [[ "$1" != "upgrade" ]]; then
	      rm -f /etc/default/$APP_NAME
	      which systemctl &>/dev/null
	      if [[ $? -eq 0 ]]; then
	         disable_systemd
	      else
	         disable_update_rcd
	      fi
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
	    if [[ "$1" = "0" ]]; then
	      rm -f /etc/default/$APP_NAME
	      disable_chkconfig
	    fi
    fi
fi

remove_traces
