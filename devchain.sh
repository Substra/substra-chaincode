#! /usr/bin/env bash

set -euo pipefail

# This script will help with chaincode dev mode setup, as explained here:
# https://hyperledger-fabric.readthedocs.io/en/latest/peer-chaincode-devmode.html

FABRIC_GIT=git@github.com:hyperledger/fabric.git
CLONE_DIR=$HOME/fabric
DATA_DIR=$HOME/fabric/devmode-data
CHAINCODE_DIR=$(pwd)

function usage() {
    local txt=(
        "Development tool to work with local chaincode"
        "Usage: ./devchain.sh [options] <command> [arguments]"
        ""
        "Command:"
        "  clean                Kill previous peer/orderer/chaincode and reset data directory."
        "  run (default)        Launch local fabric network, build and run chaincode."
        ""
        "Options:"
        "  --help, -h     Print help."
    )

    printf "%s\n" "${txt[@]}"
}

function clean() {
    echo "Cleaning previous devchain run"
    set -x
    killall orderer || true
    killall peer || true
    killall chaincode || true
    rm -rf $DATA_DIR
}

trap ctrl_c INT QUIT
function ctrl_c() {
    if [[ -n ${CHAINCODE_PID-} ]]; then
        echo "Shutting down chaincode"
        kill -s TERM $CHAINCODE_PID || true
    fi
    if [[ -n ${PEER_PID-} ]]; then
        echo "Shutting down peer"
        kill -s TERM $PEER_PID || true
    fi
    if [[ -n ${ORDERER_PID-} ]]; then
        echo "Shutting down orderer"
        kill -s TERM $ORDERER_PID || true
    fi
    rm -rf $DATA_DIR
    exit 0
}

function  run() {
    if [[ ! -d $CLONE_DIR ]]; then
        git clone $FABRIC_GIT $CLONE_DIR
    fi

    if [[ ! -d $DATA_DIR ]]; then
        mkdir -p $DATA_DIR
    fi

    cd $CLONE_DIR

    make orderer peer configtxgen

    export PATH=$CLONE_DIR/build/bin:$PATH
    export FABRIC_CFG_PATH=$CLONE_DIR/sampleconfig

    # Alter configuration to not write on /var/hyperledger as that would require priviledged execution
    ORDERER_CONF=$FABRIC_CFG_PATH/orderer.yaml
    CORE_CONF=$FABRIC_CFG_PATH/core.yaml

    sed -E -i'.orig' 's@/var/hyperledger/production@'"$DATA_DIR"'@' $CORE_CONF
    sed -E -i'.orig' 's@listenAddress: 127.0.0.1:9443@listenAddress: 127.0.0.1:9444@' $CORE_CONF
    sed -E -i'.orig' 's@/var/hyperledger/production@'"$DATA_DIR"'@' $ORDERER_CONF
    rm $CORE_CONF.orig
    rm $ORDERER_CONF.orig

    configtxgen -profile SampleDevModeSolo -channelID syschannel -outputBlock genesisblock -configPath $FABRIC_CFG_PATH -outputBlock $CLONE_DIR/sampleconfig/genesisblock

    ORDERER_GENERAL_GENESISPROFILE=SampleDevModeSolo orderer &
    ORDERER_PID=$!
    echo "Started orderer with PID: $ORDERER_PID"

    FABRIC_LOGGING_SPEC=chaincode=debug CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052 peer node start --peer-chaincodedev=true &
    PEER_PID=$!

    # Creating channel
    configtxgen -channelID ch1 -outputCreateChannelTx ch1.tx -profile SampleSingleMSPChannel -configPath $FABRIC_CFG_PATH
    peer channel create -o 127.0.0.1:7050 -c ch1 -f ch1.tx
    sleep 2

    # Join channel
    peer channel join -b ch1.block

    # Run chaincode
    echo "Building chaincode"
    cd $CHAINCODE_DIR; go build .; cd -
    echo "Running chaincode"
    CORE_CHAINCODE_LOGLEVEL=debug CORE_PEER_TLS_ENABLED=false CORE_CHAINCODE_ID_NAME=mycc:1.0 $CHAINCODE_DIR/chaincode -peer.address 127.0.0.1:7052 &
    CHAINCODE_PID=$!

    sleep 3

    # Approve chaincode
    peer lifecycle chaincode approveformyorg  -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --package-id mycc:1.0
    peer lifecycle chaincode checkcommitreadiness -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')"
    peer lifecycle chaincode commit -o 127.0.0.1:7050 --channelID ch1 --name mycc --version 1.0 --sequence 1 --init-required --signature-policy "OR ('SampleOrg.member')" --peerAddresses 127.0.0.1:7051

    # Init chaincode
    CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":[]}' --isInit

    until [[ ${input-} == "quit" ]]
    do
        cat << EOF
======== Press Ctrl-C to exit ========
Type this in another terminal to invoke chaincode:

export PATH=$CLONE_DIR/build/bin:\$PATH
export FABRIC_CFG_PATH=$CLONE_DIR/sampleconfig

CORE_PEER_ADDRESS=127.0.0.1:7051 peer chaincode invoke -o 127.0.0.1:7050 -C ch1 -n mycc -c '{"Args":["QueryNodes"]}'
EOF
        read input
    done
    ctrl_c
}

while (( $# ))
do
    case "$1" in
        --help | -h)
            usage
            exit 0
            ;;
        clean)
            clean
            exit 0
            ;;
        run)
            run
            exit 0
            ;;
    esac
done

# No command -> default to running
if [[ ${#} -eq 0 ]]; then
    run
    exit
fi
