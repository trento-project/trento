#!/bin/bash
TARGET_PUBLIC_HOST=$1
TARGET_HOST=$2
CONSUL_HOST=$3
ACTION=${4:-"deploy-agent"}
TRENTO_BIN=${TRENTO_BIN:-"trento"}
CONSUL_BIN=${CONSUL_BIN:-"consul"}
TRENTO_PATH=${TRENTO_PATH:-"/srv/trento/"}
CONSUL_PATH=${CONSUL_PATH:-"/srv/consul/"}
CONSUL_DATA_DIR=${CONSUL_DATA_DIR:-"consul-agent-data"}
TARGET_USER=${TARGET_USER:-"root"}

# Abort when any command in the script fails
set -e

# Abort if no input params
if [ $# -lt 2 ] ; then
    echo "Usage: ./deploy.sh <target-public-server-ip> <target-private-ip> <consul-ip> [deploy-agent*|deploy-web]"
    exit 1
fi

stop_process () {
    echo "Checking if process $2 is running in $1..."
    while ssh "$TARGET_USER@$TARGET_PUBLIC_HOST" pgrep -x "$2" > /dev/null
    do
        echo "Attempting to stop $2 process on $1..."
        ssh "$TARGET_USER@$TARGET_PUBLIC_HOST" killall -5 "$2"
        sleep 2
    done
}

# Stop old processes
stop_process "$TARGET_PUBLIC_HOST" "trento"
stop_process "$TARGET_PUBLIC_HOST" "consul"

# Create directory structure if it doesn't exist
ssh "$TARGET_USER@$TARGET_PUBLIC_HOST" mkdir -p "$TRENTO_PATH" || true
ssh "$TARGET_USER@$TARGET_PUBLIC_HOST" mkdir -p "$CONSUL_PATH/consul.d" || true

# Upload new binaries & examples
rsync -av ./$TRENTO_BIN "$TARGET_USER@$TARGET_PUBLIC_HOST:/$TRENTO_PATH"
rsync -av ./$CONSUL_BIN "$TARGET_USER@$TARGET_PUBLIC_HOST:/$CONSUL_PATH"

# Give them execution permission
ssh "$TARGET_USER@$TARGET_PUBLIC_HOST" chmod +x "$TRENTO_PATH/$TRENTO_BIN"
ssh "$TARGET_USER@$TARGET_PUBLIC_HOST" chmod +x "$CONSUL_PATH/$CONSUL_BIN"

if test "$TARGET_USER" != "root" ; then
    ADD_PARAMS="sudo"
fi

# Start 'em
if [ "$ACTION" = "deploy-agent" ] ; then
	ssh -t "$TARGET_USER@$TARGET_PUBLIC_HOST" -f "nohup $CONSUL_PATH/$CONSUL_BIN agent -bind=$TARGET_HOST -client=$TARGET_HOST -data-dir=$CONSUL_DATA_DIR -config-dir=$CONSUL_PATH/consul.d -retry-join=$CONSUL_HOST -ui > /dev/null 2>&1"
  ssh -t "$TARGET_USER@$TARGET_PUBLIC_HOST" -f "export CONSUL_HTTP_ADDR=$TARGET_HOST:8500 && nohup $ADD_PARAMS $TRENTO_PATH/$TRENTO_BIN agent start --consul-config-dir=$CONSUL_PATH/consul.d > /dev/null 2>&1"
elif [ "$ACTION" = "deploy-web" ] ; then
	ssh -t "$TARGET_USER@$TARGET_PUBLIC_HOST" -f "nohup $CONSUL_PATH/$CONSUL_BIN agent -server -bootstrap-expect=1 -bind=$CONSUL_HOST -data-dir=$CONSUL_DATA_DIR -client='$CONSUL_HOST $TARGET_PUBLIC_HOST' -ui > /dev/null 2>&1"
	ssh -t "$TARGET_USER@$TARGET_PUBLIC_HOST" -f "export CONSUL_HTTP_ADDR=$CONSUL_HOST:8500 && nohup $TRENTO_PATH/$TRENTO_BIN web serve > /dev/null 2>&1"
fi
