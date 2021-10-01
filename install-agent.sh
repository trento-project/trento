#!/bin/bash
# The script requires root permissions

set -e

if [ "$EUID" -ne 0 ]; then
    echo "Please run as root."
    exit
fi

function print_help() {
    cat <<END
This is a trento-agent installer. Trento is a web-based graphical user interface
that can help you deploy, provision and operate infrastructure for SAP Applications

Usage:

  sudo ./install-agent.sh --agent-bind-ip <192.168.122.10> --server-ip <192.168.122.5>

Arguments:
  --agent-bind-ip   The private address to which the trento-agent should be bound for internal communications.
                    This is an IP address that should be reachable by the other hosts, including the trento server.
  --server-ip       The trento server ip.
  --rolling         Use the factory/rolling-release version instead of the stable one.
  --use-tgz         Use the trento tar.gz file from GH releases rather than the RPM
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
    "agent-bind-ip:"
    "server-ip:"
    "rolling"
    "use-tgz"
)

readonly TRENTO_VERSION=0.4.0

opts=$(
    getopt \
        --longoptions "$(printf "%s," "${ARGUMENT_LIST[@]}")" \
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

    --rolling)
        USE_ROLLING=1
        shift 1
        ;;

    --use-tgz)
        USE_TGZ=1
        shift 1
        ;;

    *)
        break
        ;;
    esac
done

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
    echo "Warning: non-SUSE operating system, forcing --use-tgz"
    USE_TGZ=1
fi

echo "Installing trento-agent..."

function check_installer_deps() {
    if ! which unzip >/dev/null 2>&1; then
        echo "unzip is required by this script. Please install it with: zypper in -y unzip"
        exit 1
    fi
    if ! which curl >/dev/null 2>&1; then
        echo "curl is required by this script. Please install it with: zypper in -y curl"
        exit 1
    fi
}

function configure_installation() {
    if [[ -z "$AGENT_BIND_IP" ]]; then
        read -rp "Please provide a bind IP for the agent: " AGENT_BIND_IP </dev/tty
    fi
    if [[ -z "$SERVER_IP" ]]; then
        read -rp "Please provide the server IP: " SERVER_IP </dev/tty
    fi
}

function install_consul() {
    mkdir -p $CONFIG_PATH
    pushd -- "$CONSUL_PATH" >/dev/null
    curl -f -sS -O -L "https://releases.hashicorp.com/consul/$CONSUL_VERSION/consul_${CONSUL_VERSION}_linux_amd64.zip" >/dev/null
    unzip -o "consul_${CONSUL_VERSION}_linux_amd64".zip >/dev/null
    rm "consul_${CONSUL_VERSION}_linux_amd64".zip
    popd >/dev/null
}

function setup_consul() {
    echo "$CONSUL_HCL_TEMPLATE" |
        sed "s|@JOIN_ADDR@|${SERVER_IP}|g" |
        sed "s|@BIND_ADDR@|${AGENT_BIND_IP}|g" \
            >${CONFIG_PATH}/consul.hcl

    if [[ -f "/usr/lib/systemd/system/$CONSUL_SERVICE_NAME" ]]; then
        echo "  Warning: Consul systemd unit already installed. Removing..."
        systemctl stop "$CONSUL_SERVICE_NAME"
        rm "/usr/lib/systemd/system/$CONSUL_SERVICE_NAME"
    fi

    echo "$CONSUL_SERVICE_TEMPLATE" >/usr/lib/systemd/system/$CONSUL_SERVICE_NAME
    systemctl daemon-reload
}

function install_trento() {
    if [[ -n "$USE_TGZ" ]] ; then
        install_trento_tgz
    else
        install_trento_rpm
    fi
}

function install_trento_rpm() {
    if [[ -n "$USE_ROLLING" ]] ; then
        TRENTO_REPO=${TRENTO_REPO:-"https://download.opensuse.org/repositories/devel:/sap:/trento:/factory/15.3/devel:sap:trento:factory.repo"}
        TRENTO_REPO_KEY=${TRENTO_REPO_KEY:-"https://download.opensuse.org/repositories/devel:/sap:/trento:/factory/15.3/repodata/repomd.xml.key"}
    else
        TRENTO_REPO=${TRENTO_REPO:-"https://download.opensuse.org/repositories/devel:/sap:/trento/15.3/devel:sap:trento.repo"}
        TRENTO_REPO_KEY=${TRENTO_REPO_KEY:-"https://download.opensuse.org/repositories/devel:/sap:/trento/15.3/repodata/repomd.xml.key"}        
    fi

    rpm --import "${TRENTO_REPO_KEY}" >/dev/null
    path=${TRENTO_REPO%/*}/
    if zypper lr --details | cut -d'|' -f9 | grep "$path" >/dev/null 2>&1; then
        echo "* $path repository already exists. Skipping."
    else
        echo "* Adding Trento repository: $path."
        zypper ar "$TRENTO_REPO" >/dev/null
    fi
    zypper ref >/dev/null
    if which trento >/dev/null 2>&1; then
        echo "* Trento is already installed. Updating trento"
        zypper up -y trento >/dev/null
    else
        echo "* Installing trento"
        zypper in -y trento >/dev/null
    fi
}

function install_trento_tgz() {
    ARCH=$(uname -m | sed "s~x86_64~amd64~" | sed "s~aarch64~arm64~" )
    local bin_dir=${BIN_DIR:-"/usr/bin"}
    local sysd_dir=${SYSD_DIR:-"/usr/lib/systemd/system"}
    
    if [[ -n "$USE_ROLLING" ]] ; then
        TRENTO_TGZ_URL=https://github.com/trento-project/trento/releases/download/rolling/trento-${ARCH}.tgz
    else
        TRENTO_TGZ_URL=https://github.com/trento-project/trento/releases/download/${TRENTO_VERSION}/trento-${ARCH}.tgz
    fi
  
    curl -f -sS -O -L "${TRENTO_TGZ_URL}" >/dev/null
    tar -zxf trento-${ARCH}.tgz

    mv trento ${bin_dir}/trento
    mv packaging/systemd/trento-agent.service ${sysd_dir}/trento-agent.service
    systemctl daemon-reload
    systemctl enable --now trento-agent.service
    rm trento-${ARCH}.tgz
}

check_installer_deps
configure_installation
install_consul
setup_consul
install_trento

echo -e "\e[92mDone.\e[97m"
echo -e "Now you can start trento-agent with: \033[1msystemctl start trento-agent\033[0m"
echo -e "Please make sure the \033[1mserver\033[0m is running before starting the agent."
