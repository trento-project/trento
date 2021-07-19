#!/bin/bash
# The script requires root permissions

set -e

if [ "$EUID" -ne 0 ]; then
    echo "Please run as root."
    exit
fi

function print_help() {
    cat << END
This is a trento-agent installer. Trento is a web-based graphical user interface
that can help you deploy, provision and operate infrastructure for SAP Applications

Usage:

  sudo ./install.sh --agent-bind-ip <192.168.122.10> --server-ip <192.168.122.5>

Arguments:
  --agent-bind-ip   The private address to which the trento-agent should be bound for internal communications.
                    This is an IP address that should be reachable by the other nodes, including the trento server.
  --server-ip       The trento server ip.
  --help            Print this help.
END
}

case "$1" in
--help)
    print_help
    exit 0
    ;;
esac

ARGUMENT_LIST=(
    "agent-bind-ip"
    "server-ip"
)

opts=$(
    getopt \
        --longoptions "$(printf "%s:," "${ARGUMENT_LIST[@]}")" \
        --name "$(basename "$0")" \
        --options "" \
        -- "$@"
)

eval set "--$opts"

while [[ $# -gt 0 ]]; do
    case "$1" in
    --agent-bind-ip)
        AGENT_BIND_IP=$2
        shift 2
        ;;

    --server-ip)
        SERVER_IP=$2
        shift 2
        ;;

    *)
        break
        ;;
    esac
done

if [ -z "$AGENT_BIND_IP" ]; then
    read -rp "Please provide a bind IP for the agent: " AGENT_BIND_IP < /dev/tty
fi
if [ -z "$SERVER_IP" ]; then
    read -rp "Please provide the server IP: " SERVER_IP < /dev/tty
fi
if [ -z "$NODE_NAME" ]; then
    NODE_NAME="$HOSTNAME"
fi

TRENTO_REPO_KEY=${TRENTO_REPO_KEY:-"https://download.opensuse.org/repositories/devel:/sap:/trento/15.3/repodata/repomd.xml.key"}
TRENTO_REPO=${TRENTO_REPO:-"https://download.opensuse.org/repositories/devel:/sap:/trento/15.3/devel:sap:trento.repo"}

CONSUL_VERSION=1.9.6
CONSUL_PATH=/srv/consul
CONFIG_PATH="$CONSUL_PATH/consul.d"
CONSUL_HCL_TEMPLATE='data_dir = "/srv/consul/data/"
log_level = "DEBUG"
datacenter = "dc1"
ui = true
bind_addr = "@BIND_ADDR@"
client_addr = "0.0.0.0"
retry_join = ["@JOIN_ADDR@"]'

CONSUL_SERVICE_NAME="consul.service"
CONSUL_SERVICE_TEMPLATE='[Unit]
Description="HashiCorp Consul - A service mesh solution"
Documentation=https://www.consul.io/
Requires=network-online.target
After=network-online.target
ConditionFileNotEmpty=/srv/consul/consul.d/consul.hcl
PartOf=trento-agent.service

[Service]
ExecStart=/srv/consul/consul agent -config-dir=/srv/consul/consul.d
ExecReload=/bin/kill --signal HUP $MAINPID
KillMode=process
Restart=on-failure
RestartSec=5
Type=notify


[Install]
WantedBy=multi-user.target'

. /etc/os-release
if [[ ! $PRETTY_NAME =~ "SUSE" ]]; then
    echo "Operating system is not supported. Exiting."
    exit 1
fi

echo "Installing trento-agent..."

function install_deps() {
    echo "* Installing dependencies... "
    if ! which unzip > /dev/null 2>&1; then
        echo "* Installing unzip"
        zypper in -y unzip > /dev/null
    fi
    if ! which curl > /dev/null 2>&1; then
        echo "* Installing curl"
        zypper in -y curl > /dev/null
    fi
}

function install_consul() {
    mkdir -p $CONFIG_PATH
    pushd -- "$CONSUL_PATH" > /dev/null
    curl -f -sS -O -L "https://releases.hashicorp.com/consul/$CONSUL_VERSION/consul_${CONSUL_VERSION}_linux_amd64.zip" > /dev/null
    unzip -o "consul_${CONSUL_VERSION}_linux_amd64".zip > /dev/null
    rm "consul_${CONSUL_VERSION}_linux_amd64".zip
    popd > /dev/null
}

function setup_consul() {
    echo "$CONSUL_HCL_TEMPLATE" |
        sed "s|@JOIN_ADDR@|${SERVER_IP}|g" |
        sed "s|@BIND_ADDR@|${AGENT_BIND_IP}|g" |
        sed "s|@NODE_NAME@|${NODE_NAME}|g" \
            >${CONFIG_PATH}/consul.hcl

    if [ -f "/usr/lib/systemd/system/$CONSUL_SERVICE_NAME" ]; then
        echo "  Warning: Consul systemd unit already installed. Removing..."
        systemctl stop "$CONSUL_SERVICE_NAME"
        rm "/usr/lib/systemd/system/$CONSUL_SERVICE_NAME"
    fi

    echo "$CONSUL_SERVICE_TEMPLATE" >/usr/lib/systemd/system/$CONSUL_SERVICE_NAME
    systemctl daemon-reload
}

function install_trento() {
    rpm --import "${TRENTO_REPO_KEY}" > /dev/null
    path=${TRENTO_REPO%/*}/
    if zypper lr --details | cut -d'|' -f9 | grep "$path" > /dev/null 2>&1; then
        echo "* $path repository already exists. Skipping."
    else
        zypper ar "$TRENTO_REPO" > /dev/null
    fi
    zypper ref > /dev/null
    if which trento > /dev/null 2>&1; then
        echo "* Trento is already installed. Updating trento"
        zypper up -y trento > /dev/null
    else
        echo "* Installing trento"
        zypper in -y trento > /dev/null
    fi
}

install_consul
setup_consul
install_trento

echo -e "\e[92mDone.\e[97m"
echo -e "You can now start trento-agent with: \033[1msystemctl start trento-agent\033[0m"
