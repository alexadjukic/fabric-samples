#!/bin/bash

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config

. scripts/utils.sh
. scripts/envVar.sh

# invokeInit ORG PEER (ORG PEER...)
function invokeInit() {
  parsePeerConnectionParameters $@

  # set -x
  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c "{\"function\":\"InitLedger\",\"Args\":[\"$(date +%a\ %b\ %e\ %T\ %Z\ %Y)\"]}"
}

# getAll ORG PEER ASSET_NAME
function getAll() {
  ORG=$1
  PEER=$2
  ASSET_NAME=$3
  setGlobals $ORG $PEER
  peer chaincode query -C $CHANNEL_NAME -n $CC_NAME -c "{\"function\": \"GetAll${ASSET_NAME}\", \"Args\":[]}" | jq
}


FUNC=$1
CHANNEL_NAME=${2:-"mychannel"}
CC_NAME=${3:-"basic"}

if [ $FUNC == "InitLedger" ]; then
  infoln "Invoking InitLedger function"
  invokeInit 1 0 2 0 3 0
elif [ $FUNC == "GetAllUsers" ]; then
  getAll 3 0 "Users"
elif [ $FUNC == "GetAllProducts" ]; then
  getAll 3 0 "Products"
elif [ $FUNC == "GetAllMerchants" ]; then
  getAll 3 0 "Merchants"
fi

