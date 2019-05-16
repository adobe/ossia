#!/bin/bash

APP_NAME=ossia
BIN_DIR=$APP_NAME/bin
DATA_DIR=/opt/$APP_NAME
LOG_DIR=/opt/$APP_NAME/log
SCRIPT_DIR=/opt/$APP_NAME/scripts
LOGROTATE_DIR=/etc/logrotate.d

function install_init {
    cp -f $SCRIPT_DIR/init.sh /etc/init.d/$APP_NAME
    chmod +x /etc/init.d/$APP_NAME
}

function install_upstart {
    cp -f $SCRIPT_DIR/$APP_NAME.conf /etc/init/
}

function install_systemd {
    cp -f $SCRIPT_DIR/$APP_NAME.service /lib/systemd/system/$APP_NAME.service
    systemctl enable $APP_NAME
}

function install_update_rcd {
    update-rc.d $APP_NAME defaults
}

function install_chkconfig {
    chkconfig --add $APP_NAME
}

id $APP_NAME &>/dev/null
if [[ $? -ne 0 ]]; then
    useradd --system -U -M $APP_NAME -s /bin/false -d $DATA_DIR
fi

chown -R -L $APP_NAME:$APP_NAME $DATA_DIR
chown -R -L $APP_NAME:$APP_NAME $LOG_DIR

# Add defaults file, if it doesn't exist
if [[ ! -f /etc/default/$APP_NAME ]]; then
    touch /etc/default/$APP_NAME
fi

# Remove legacy symlink, if it exists
if [[ -L /etc/init.d/$APP_NAME ]]; then
    rm -f /etc/init.d/$APP_NAME
fi

# Distribution-specific logic
if [[ -f /etc/redhat-release ]]; then
    # RHEL-variant logic
    which systemctl &>/dev/null
    if [[ $? -eq 0 ]]; then
       install_systemd
    else
       install_init
       install_chkconfig
    fi
elif [[ -f /etc/debian_version ]]; then
    # Debian/Ubuntu logic
    which systemctl &>/dev/null
    if [[ $? -eq 0 ]]; then
      install_systemd
    else
      # Assuming sysv
      install_init
      install_update_rcd
      #install_upstart
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ $ID = "amzn" ]]; then
       # Amazon Linux logic
       install_init
       install_chkconfig
    fi
fi
