package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.InitLedger(transactionContext, "Tue Dec  1 01:34:07 CET 2037")
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = assetTransfer.InitLedger(transactionContext, "Tue Dec  1 01:34:07 CET 2037")
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

func TestCreateMerchant(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.CreateMerchant(transactionContext, "", "", 0, 1000)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.CreateMerchant(transactionContext, "21", "", 0, 1000)
	require.EqualError(t, err, "the merchant 21 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.CreateMerchant(transactionContext, "21", "", 0, 1000)
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestAddProductToMerchant(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	assetTransfer := chaincode.SmartContract{}

	merchant := chaincode.Merchant{
		ID:       "21",
		Pib:      "1234",
		Type:     chaincode.Wholesale,
		Products: []string{"11", "12"},
		Bills:    []string{},
		Balance:  1000,
	}

	product := chaincode.Product{
		ID:         "13",
		Name:       "Toilet paper",
		ExpiryDate: "",
		Price:      6,
		Amount:     20,
	}

	merchantJSON, err := json.Marshal(merchant)
	require.NoError(t, err)

	productJSON, err := json.Marshal(product)
	require.NoError(t, err)

	chaincodeStub.GetStateReturnsOnCall(0, merchantJSON, nil)
	chaincodeStub.GetStateReturnsOnCall(1, productJSON, nil)
	err = assetTransfer.AddProductToMerchant(transactionContext, "21", "13")
	require.NoError(t, err)

	chaincodeStub.GetStateReturnsOnCall(2, merchantJSON, nil)
	chaincodeStub.GetStateReturnsOnCall(3, productJSON, nil)
	err = assetTransfer.AddProductToMerchant(transactionContext, "21", "11")
	require.EqualError(t, err, "merchant 21 already has product 11")

	chaincodeStub.GetStateReturnsOnCall(4, nil, nil)
	err = assetTransfer.AddProductToMerchant(transactionContext, "21", "11")
	require.EqualError(t, err, "the merchant 21 does not exist")

	chaincodeStub.GetStateReturnsOnCall(5, merchantJSON, nil)
	chaincodeStub.GetStateReturnsOnCall(6, nil, nil)
	err = assetTransfer.AddProductToMerchant(transactionContext, "21", "11")
	require.EqualError(t, err, "the product 11 does not exist")
}

func TestCreateUsers(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}

	users := []chaincode.User{
		{
			Balance: 1000,
			Bills:   []string{},
			Email:   "markomarkovic@email.com",
			ID:      "01",
			Name:    "Marko",
			Surname: "Markovic",
		},
		{
			Balance: 1000,
			Bills:   []string{},
			Email:   "petarpetrovic@email.com",
			ID:      "02",
			Name:    "Petar",
			Surname: "Petrovic",
		},
	}

	chaincodeStub.GetStateReturns(nil, nil)
	err := assetTransfer.CreateUsers(transactionContext, users)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.CreateUsers(transactionContext, users)
	require.EqualError(t, err, "the user 01 already exists")
}

func TestCreateAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.CreateAsset(transactionContext, "", "", 0, "", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = assetTransfer.CreateAsset(transactionContext, "asset1", "", 0, "", 0)
	require.EqualError(t, err, "the asset asset1 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.CreateAsset(transactionContext, "asset1", "", 0, "", 0)
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestReadAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.Asset{ID: "asset1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	asset, err := assetTransfer.ReadAsset(transactionContext, "")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.ReadAsset(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	asset, err = assetTransfer.ReadAsset(transactionContext, "asset1")
	require.EqualError(t, err, "the asset asset1 does not exist")
	require.Nil(t, asset)
}

func TestUpdateAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.Asset{ID: "asset1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.UpdateAsset(transactionContext, "", "", 0, "", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.UpdateAsset(transactionContext, "asset1", "", 0, "", 0)
	require.EqualError(t, err, "the asset asset1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.UpdateAsset(transactionContext, "asset1", "", 0, "", 0)
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestDeleteAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	asset := &chaincode.Asset{ID: "asset1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.DeleteAsset(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = assetTransfer.DeleteAsset(transactionContext, "asset1")
	require.EqualError(t, err, "the asset asset1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.DeleteAsset(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestTransferAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	asset := &chaincode.Asset{ID: "asset1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	err = assetTransfer.TransferAsset(transactionContext, "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = assetTransfer.TransferAsset(transactionContext, "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

// func TestGetAllAssets(t *testing.T) {
// 	asset := &chaincode.Asset{ID: "asset1"}
// 	bytes, err := json.Marshal(asset)
// 	require.NoError(t, err)

// 	iterator := &mocks.StateQueryIterator{}
// 	iterator.HasNextReturnsOnCall(0, true)
// 	iterator.HasNextReturnsOnCall(1, false)
// 	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

// 	chaincodeStub := &mocks.ChaincodeStub{}
// 	transactionContext := &mocks.TransactionContext{}
// 	transactionContext.GetStubReturns(chaincodeStub)

// 	chaincodeStub.GetStateByRangeReturns(iterator, nil)
// 	assetTransfer := &chaincode.SmartContract{}
// 	assets, err := assetTransfer.GetAllAssets(transactionContext)
// 	require.NoError(t, err)
// 	require.Equal(t, []*chaincode.Asset{asset}, assets)

// 	iterator.HasNextReturns(true)
// 	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
// 	assets, err = assetTransfer.GetAllAssets(transactionContext)
// 	require.EqualError(t, err, "failed retrieving next item")
// 	require.Nil(t, assets)

// 	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
// 	assets, err = assetTransfer.GetAllAssets(transactionContext)
// 	require.EqualError(t, err, "failed retrieving all assets")
// 	require.Nil(t, assets)
// }
