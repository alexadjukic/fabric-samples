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

# createMerchant ORG PEER (ORG PEER...)
function createMerchant() {
  parsePeerConnectionParameters $@

  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function":"CreateMerchant","Args":["23", "789", "Supermarket", "1000"]}'
}

function addProductToMerchant() {
  parsePeerConnectionParameters $@

  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function":"AddProductToMerchant","Args":["21", "13"]}'
}

# createUsers ORG PEER (ORG PEER...)
function createUsers() {
  parsePeerConnectionParameters $@

  JSON_PAYLOAD=$(cat <<EOF
[
  {
    "ID": "03",
    "Name": "John",
    "Surname": "Doe",
    "Email": "john@example.com",
    "Balance": 100,
    "Bills": []
  },
  {
    "ID": "04",
    "Name": "Jane",
    "Surname": "Smith",
    "Email": "jane@example.com",
    "Balance": 200,
    "Bills": []
  }
]
EOF
  )

  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c "$(jq -nc --argjson users "$JSON_PAYLOAD" '{"function": "CreateUsers", "Args": [($users | tostring)]}')"
}

# buyProduct ORG PEER (ORG PEER...)
function buyProduct() {
  parsePeerConnectionParameters $@

  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c "{\"function\": \"BuyProduct\", \"Args\": [\"01\", \"11\", \"$(date +%a\ %b\ %e\ %T\ %Z\ %Y)\"]}"
}

# addFunds ORG PEER (ORG PEER...)
function addFunds() {
  parsePeerConnectionParameters $@

  peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function": "AddFunds", "Args": ["01", "5"]}'
}

# findProduct ORG PEER ID NAME MERCHANT_TYPE PRICE
function findProduct() {
  ORG=$1
  PEER=$2
  ID=$3
  NAME=$4
  MERCHANT_TYPE=$5
  PRICE=$6
  setGlobals $ORG $PEER
  peer chaincode query -C $CHANNEL_NAME -n $CC_NAME -c "{\"function\": \"FindProduct\", \"Args\":[\"${ID}\", \"${NAME}\", \"${MERCHANT_TYPE}\", \"${PRICE}\"]}" | jq
}


FUNC=$1
CHANNEL_NAME=${2:-"mychannel"}
CC_NAME=${3:-"basic"}

if [ $FUNC == "InitLedger" ]; then
  infoln "Invoking InitLedger function"
  invokeInit 1 0 2 0 3 0
elif [ $FUNC == "GetAllUsers" ]; then
  infoln "Invoking GetAllUsers function"
  getAll 3 0 "Users"
elif [ $FUNC == "GetAllProducts" ]; then
  infoln "Invoking GetAllProducts function"
  getAll 3 0 "Products"
elif [ $FUNC == "GetAllMerchants" ]; then
  infoln "Invoking GetAllMerchants function"
  getAll 3 0 "Merchants"
elif [ $FUNC == "CreateMerchant" ]; then
  infoln "Invoking CreateMerchant function with parameters"
  infoln "  id: 23"
  infoln "  pib: 789"
  infoln "  type: Supermarket"
  infoln "  balance: 1000"
  createMerchant 1 0 2 0 3 0
elif [ $FUNC == "AddProductToMerchant" ]; then
  infoln "Invoking AddProductToMerchant function with parameters"
  infoln "  merchantId: 21"
  infoln "  productId: 13"
  addProductToMerchant 1 0 2 0 3 0
elif [ $FUNC == "CreateUsers" ]; then
  infoln "Invoking CreateUsers function with 2 users"
  createUsers 1 0 2 0 3 0
elif [ $FUNC == "BuyProduct" ]; then
  infoln "Invoking BuyProduct function with arguments"
  infoln "  userId: 01"
  infoln "  productId: 11"
  buyProduct 1 0 2 0 3 0
elif [ $FUNC == "AddFunds" ]; then
  infoln "Invoking AddFunds function with parameters"
  infoln "  id: 01"
  infoln "  amount: 5"
  addFunds 1 0 2 0 3 0
elif [ $FUNC == "FindProduct" ]; then
  infoln "Invoking FindProduct function with parameters"
  infoln "  id: 11"
  infoln "  name: Cheese"
  infoln "  mtype: Wholesale"
  infoln "  price: 5"
  findProduct 3 0 "11" "Cheese" "Wholesale" "5"
fi

