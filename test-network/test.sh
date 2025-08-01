./network.sh down
./network.sh up createChannel -ca -s couchdb
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
./invokeCC.sh InitLedger
