#!/bin/bash

set -eu

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

echo ""
echo "Provisioning..."

if [ $# -lt 7 ]; then
  echo "Usage: ./provision.sh <github-user> <github-repo> <actions-runner-token> <actions-runner-name> <consul-bind-ip> <consul-server-ip> [agent|web]"
  exit 1
fi

ACTIONS_RUNNER_VERSION=2.278.0
ACTIONS_RUNNER_USER=github-runner
ACTIONS_RUNNER_HOME=/srv/github-runner
ACTIONS_RUNNER_REPO_OWNER=$1
ACTIONS_RUNNER_REPO_NAME=$2
ACTIONS_RUNNER_TOKEN=$3
ACTIONS_RUNNER_NAME=$4
ACTIONS_RUNNER_PATH="$ACTIONS_RUNNER_HOME/actions-runner"
ACTIONS_RUNNER_REPO_URL="https://github.com/$ACTIONS_RUNNER_REPO_OWNER/$ACTIONS_RUNNER_REPO_NAME"
ACTIONS_RUNNER_SYSTEMD_UNIT="actions.runner.$ACTIONS_RUNNER_REPO_OWNER-$ACTIONS_RUNNER_REPO_NAME.$ACTIONS_RUNNER_NAME.service"

CONSUL_VERSION=1.9.6
CONSUL_BIND_IP=$5
CONSUL_SERVER_IP=$6
CONSUL_USER=consul
CONSUL_HOME=/srv/consul
CONSUL_CONFIG_PATH="$CONSUL_HOME/consul.d"

ROLE=$7
if [ "$ROLE" = "agent" ]; then
  CONSUL_HCL="consul-client.hcl"
elif [ "$ROLE" = "web" ]; then
  CONSUL_HCL="consul-server.hcl"
else
  echo "Please specify a valid role"
  exit
fi

CONSUL_HCL_TEMPLATE_URL="https://raw.githubusercontent.com/$ACTIONS_RUNNER_REPO_OWNER/$ACTIONS_RUNNER_REPO_NAME/main/hack/$CONSUL_HCL.template"

CONSUL_SYSTEMD_UNIT="consul.service"
CONSUL_SYSTEMD_UNIT_URL="https://raw.githubusercontent.com/$ACTIONS_RUNNER_REPO_OWNER/$ACTIONS_RUNNER_REPO_NAME/main/hack/$CONSUL_SYSTEMD_UNIT"

TRENTO_PATH=/srv/trento
TRENTO_SYSTEMD_UNIT="trento-$ROLE.service"
TRENTO_SYSTEMD_UNIT_URL="https://raw.githubusercontent.com/$ACTIONS_RUNNER_REPO_OWNER/$ACTIONS_RUNNER_REPO_NAME/main/hack/$TRENTO_SYSTEMD_UNIT"

create_user() {
  local user=$1
  local home=$2

  echo ""
  echo "* Creating user: $user"

  if id "$user" &>/dev/null; then
    echo "  Warning: user $user already exists. Skipping..."
    return
  fi

  useradd --system -d "$home" "$user"
  mkdir -p "$home"
  chown "$user" "$home"
}

install_actions_runner() {
  echo ""
  echo "* Installing GitHub Actions Runner"

  if [ -f "$ACTIONS_RUNNER_PATH/.runner" ]; then
    echo "  Warning: Actions Runner already installed and configured. Skipping..."
    return
  fi

  mkdir -p $ACTIONS_RUNNER_PATH
  pushd -- "$ACTIONS_RUNNER_PATH" >/dev/null
  curl -f -sS -O -L "https://github.com/actions/runner/releases/download/v${ACTIONS_RUNNER_VERSION}/actions-runner-linux-x64-${ACTIONS_RUNNER_VERSION}.tar.gz" >/dev/null
  tar xfz "actions-runner-linux-x64-${ACTIONS_RUNNER_VERSION}.tar.gz"
  rm "actions-runner-linux-x64-${ACTIONS_RUNNER_VERSION}.tar.gz"
  chown -R $ACTIONS_RUNNER_USER $ACTIONS_RUNNER_PATH

  echo "  Installing dependencies"
  ./bin/installdependencies.sh >/dev/null

  echo "  Configuring"
  su -m $ACTIONS_RUNNER_USER -c "./config.sh --token $ACTIONS_RUNNER_TOKEN --unattended --url $ACTIONS_RUNNER_REPO_URL --name $ACTIONS_RUNNER_NAME --labels $ACTIONS_RUNNER_NAME >/dev/null"

  if [ -f "/etc/systemd/system/$ACTIONS_RUNNER_SYSTEMD_UNIT" ]; then
    echo "  Warning: Systemd unit already installed. Removing..."
    systemctl stop "$ACTIONS_RUNNER_SYSTEMD_UNIT"
    rm "/etc/systemd/system/$ACTIONS_RUNNER_SYSTEMD_UNIT"
  fi

  echo "  Installing systemd unit"
  ./svc.sh install >/dev/null
  systemctl enable --now "$ACTIONS_RUNNER_SYSTEMD_UNIT"

  popd >/dev/null
}

install_consul() {
  echo ""
  echo "* Installing Consul"

  mkdir -p $CONSUL_CONFIG_PATH
  pushd -- "$CONSUL_HOME" >/dev/null
  curl -f -sS -O -L "https://releases.hashicorp.com/consul/$CONSUL_VERSION/consul_${CONSUL_VERSION}_linux_amd64.zip" >/dev/null
  unzip -o "consul_${CONSUL_VERSION}_linux_amd64".zip >/dev/null
  rm "consul_${CONSUL_VERSION}_linux_amd64".zip
  chown -R $CONSUL_USER $CONSUL_HOME
  popd >/dev/null
}

setup_consul() {
  echo ""
  echo "* Setting up Consul"

  echo "  Creating configuration"
  pushd -- $CONSUL_CONFIG_PATH >/dev/null
  curl -f -sS -O -L $CONSUL_HCL_TEMPLATE_URL >/dev/null
  cat $CONSUL_HCL.template | sed "s|@JOIN_ADDR@|${CONSUL_SERVER_IP}|g" | sed "s|@BIND_ADDR@|${CONSUL_BIND_IP}|g" >consul.hcl
  rm $CONSUL_HCL.template
  popd >/dev/null

  if [ -f "/etc/systemd/system/$CONSUL_SYSTEMD_UNIT" ]; then
    echo "  Warning: Systemd unit already installed. Removing..."
    systemctl stop "$CONSUL_SYSTEMD_UNIT"
    rm "/etc/systemd/system/$CONSUL_SYSTEMD_UNIT"
  fi

  echo "  Installing systemd unit"
  curl -f -sS -L $CONSUL_SYSTEMD_UNIT_URL -o /tmp/$CONSUL_SYSTEMD_UNIT >/dev/null

  if [ "$ROLE" = "web" ]; then
    sed -i "s|Type=notify|Type=simple|g" /tmp/$CONSUL_SYSTEMD_UNIT
  fi

  mv /tmp/$CONSUL_SYSTEMD_UNIT /etc/systemd/system/
  systemctl daemon-reload
  systemctl enable --now $CONSUL_SYSTEMD_UNIT

  echo "  Adding sudoers entries"
  echo "$ACTIONS_RUNNER_USER ALL=(ALL) NOPASSWD: /bin/systemctl start consul" >/etc/sudoers.d/github-runner-consul
  echo "$ACTIONS_RUNNER_USER ALL=(ALL) NOPASSWD: /bin/systemctl stop consul" >>/etc/sudoers.d/github-runner-consul
  echo "$ACTIONS_RUNNER_USER ALL=(ALL) NOPASSWD: /bin/rm -rf /srv/consul/data" >>/etc/sudoers.d/github-runner-consul
}

setup_trento() {
  echo ""
  echo "* Setting up Trento"

  echo "  Creating Trento directory"
  mkdir -p $TRENTO_PATH
  chown -R $ACTIONS_RUNNER_USER $TRENTO_PATH

  if [ -f "/etc/systemd/system/$TRENTO_SYSTEMD_UNIT" ]; then
    echo "  Warning: Systemd unit already installed. Removing..."
    systemctl stop "$TRENTO_SYSTEMD_UNIT"
    rm "/etc/systemd/system/$TRENTO_SYSTEMD_UNIT"
  fi

  echo "  Installing systemd unit"
  curl -f -sS -L "$TRENTO_SYSTEMD_UNIT_URL" -o /tmp/"$TRENTO_SYSTEMD_UNIT" >/dev/null
  mv /tmp/"$TRENTO_SYSTEMD_UNIT" /etc/systemd/system/
  systemctl daemon-reload
  systemctl enable "$TRENTO_SYSTEMD_UNIT"

  echo "  Adding sudoers entries"
  echo "$ACTIONS_RUNNER_USER ALL=(ALL) NOPASSWD: /bin/systemctl start trento-$ROLE" >/etc/sudoers.d/github-runner-trento
  echo "$ACTIONS_RUNNER_USER ALL=(ALL) NOPASSWD: /bin/systemctl stop trento-$ROLE" >>/etc/sudoers.d/github-runner-trento
}

create_user $ACTIONS_RUNNER_USER $ACTIONS_RUNNER_HOME
install_actions_runner
create_user $CONSUL_USER $CONSUL_HOME
install_consul
setup_consul
setup_trento
echo ""
echo "Done!"
