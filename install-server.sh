#!/bin/bash

set -e

readonly ARGS=( "$@" )
readonly PROGNAME="./install-server.sh"
TRENTO_VERSION="0.6.0"

usage() {
    cat <<- EOF
    usage: $PROGNAME options

    Install Trento Server

    OPTIONS:
        -p --private-key  <file>         re-authorized private SSH key used by the runner to connect to the hosts
        -r --rolling                     use the rolling-release version instead of the stable one
        -h --help                        show this help

    Example:
       $PROGNAME --private-key ./id_rsa_runner
EOF
}

cmdline() {
    local arg=
    local private_key_absolute_path=
    for arg
    do
        local delim=""
        case "$arg" in
            --private-key)  args="${args}-p ";;
            --rolling)      args="${args}-r ";;
            --help)         args="${args}-h ";;

            # pass through anything else
            *) [[ "${arg:0:1}" == "-" ]] || delim="\""
            args="${args}${delim}${arg}${delim} ";;
        esac
    done

    eval set -- "$args"

    while getopts "p:rh" OPTION
    do
        case $OPTION in
            h)
                usage
                exit 0
            ;;
            p)
                readonly PRIVATE_KEY=$OPTARG
            ;;
            r)
                readonly ROLLING=true
            ;;
            *)
                usage
                exit 0
            ;;
        esac
    done

    if [[ -z "$PRIVATE_KEY" ]]; then
        read -rp "Please provide the path of the runner private key: " PRIVATE_KEY </dev/tty
    fi

    # Replace tilde with the current home:
    PRIVATE_KEY="${PRIVATE_KEY/#\~/$HOME}"
    private_key_absolute_path=$(realpath -q -e "$PRIVATE_KEY" || {
         echo "Path '${PRIVATE_KEY}' to private SSH key does not exist, please try again."
         exit 1
    })
    PRIVATE_KEY="$private_key_absolute_path"

    if [[ "$ROLLING" == "true" ]]; then
        TRENTO_VERSION="rolling"
    fi

    return 0
}

check_requirements() {
    local firewalld_status
    firewalld_status="$(systemctl show -p ActiveState firewalld | cut -d'=' -f2)"
    if [ "${firewalld_status}" = "active" ]; then
        echo "firewalld must be turned off to run K3s, please disable it and try again."
        exit 1
    fi
    if ! which curl >/dev/null 2>&1; then
        echo "curl is required by this script, please install it and try again."
        exit 1
    fi
    if ! which unzip >/dev/null 2>&1; then
        echo "unzip is required by this script, please install it and try again."
        exit 1
    fi
    if grep -q "Y" /sys/module/apparmor/parameters/enabled; then
        if ! command -v /sbin/apparmor_parser >/dev/null 1>&1; then
            echo "apparmor_parser is required by k3s when using AppArmor, please install it and try again."
            exit 1
        fi
    fi
}

install_k3s() {
    echo "Installing K3s..."
    curl -sfL https://get.k3s.io | sh >/dev/null
    mkdir -p ~/.kube
    sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
    sudo chown "$USER": ~/.kube/config
    unset KUBECONFIG
}

install_helm() {
    echo "Installing Helm..."
    curl -s https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash >/dev/null
}

update_helm_dependencies() {
    echo "Updating Helm dependencies..."
    helm repo add hashicorp https://helm.releases.hashicorp.com >/dev/null
    helm repo add bitnami https://charts.bitnami.com/bitnami >/dev/null
    helm repo update >/dev/null
}

install_trento_server_chart() {
    local repo_owner=${TRENTO_REPO_OWNER:-"trento-project"}
    local runner_image=${TRENTO_RUNNER_IMAGE:-"ghcr.io/$repo_owner/trento-runner"}
    local web_image=${TRENTO_WEB_IMAGE:-"ghcr.io/$repo_owner/trento-web"}
    local private_key=${PRIVATE_KEY:-"./id_rsa_runner"}
    local trento_source_zip="${TRENTO_VERSION}"
    local trento_packages_url="https://github.com/${repo_owner}/trento/archive/refs/tags"

    echo "Installing trento-server chart..."
    pushd -- /tmp >/dev/null
    rm -rf trento-"${trento_source_zip}"
    rm -f ${trento_source_zip}.zip
    curl -f -sS -O -L "${trento_packages_url}/${trento_source_zip}.zip" >/dev/null
    unzip -o "${trento_source_zip}.zip" >/dev/null
    popd >/dev/null

    pushd -- /tmp/trento-"${trento_source_zip}"/packaging/helm/trento-server >/dev/null
    helm dep update >/dev/null
    helm upgrade --install trento-server . \
    --set-file trento-runner.privateKey="${private_key}" \
    --set trento-web.image.tag="${TRENTO_VERSION}" \
    --set trento-runner.image.tag="${TRENTO_VERSION}" \
    --set trento-runner.image.repository="${runner_image}" \
    --set trento-web.image.repository="${web_image}"

    popd >/dev/null
}

main() {
    cmdline "${ARGS[@]}"
    echo "Installing trento-server on k3s..."
    check_requirements
    install_k3s
    install_helm
    update_helm_dependencies
    install_trento_server_chart
}
main
