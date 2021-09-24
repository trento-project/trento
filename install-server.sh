#!/bin/bash

set -e

readonly ARGS=( "$@" )
readonly PROGNAME="./install-server.sh"

usage() {
    cat <<- EOF
    usage: $PROGNAME options

    Install Trento Server

    OPTIONS:
       -p --private-key         pre-authorized private SSH key used by the runner to connect to the hosts
       -h --help                show this help


    Example:
       $PROGNAME --private-key ./id_rsa_runner
EOF
}

cmdline() {
    local arg=
    for arg
    do
        local delim=""
        case "$arg" in
            --private-key)  args="${args}-p ";;
            --help)         args="${args}-h ";;
            
            # pass through anything else
            *) [[ "${arg:0:1}" == "-" ]] || delim="\""
            args="${args}${delim}${arg}${delim} ";;
        esac
    done
    
    eval set -- "$args"
    
    while getopts "p:h" OPTION
    do
        case $OPTION in
            h)
                usage
                exit 0
            ;;
            p)
                readonly PRIVATE_KEY=$OPTARG
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

    return 0
}

check_deps() {
    if ! which curl >/dev/null 2>&1; then
        echo "curl is required by this script, please install it and try again."
        exit 1
    fi
    if ! which unzip >/dev/null 2>&1; then
        echo "unzip is required by this script, please install it and try again."
        exit 1
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
    # FIXME: why does /usr/local/bin vanish from PATH?
    PATH=$PATH:/usr/local/bin
    curl -s https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash >/dev/null
}

update_helm_dependencies() {
    echo "Updating Helm dependencies..."
    helm repo add hashicorp https://helm.releases.hashicorp.com >/dev/null
    helm repo update >/dev/null
}


install_trento_server_chart() {
    local repo_owner=${TRENTO_REPO_OWNER:-"trento-project"}
    local repo_branch=${TRENTO_REPO_BRANCH:-"main"}
    local private_key=${PRIVATE_KEY:-"./id_rsa_runner"}
    local image_tag=${IMAGE_TAG:-""}
    
    echo "Installing trento-server chart..."
    pushd -- /tmp >/dev/null
    curl -f -sS -O -L https://github.com/"${repo_owner}"/trento/archive/refs/heads/"${repo_branch}".zip>/dev/null
    unzip -o "${repo_branch}".zip >/dev/null
    rm "${repo_branch}".zip
    popd >/dev/null
    
    pushd -- /tmp/trento-"${repo_branch}"/packaging/helm/trento-server >/dev/null 
    helm dep update >/dev/null
    helm upgrade --install trento-server . --set-file trento-runner.privateKey="${private_key}" --set trento-web.image.tag="${image_tag}" --set trento-runner.image.tag="${image_tag}"
    rm -rf /tmp/trento-"${repo_branch}"
    popd >/dev/null
}

main() {
    cmdline "${ARGS[@]}"
    echo "Installing trento-server on k3s..."
    install_k3s
    install_helm
    update_helm_dependencies
    install_trento_server_chart
}
main