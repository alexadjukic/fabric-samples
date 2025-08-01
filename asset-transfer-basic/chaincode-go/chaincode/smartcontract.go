package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

type Merchant struct {
	Balance  int          `json:"Balance"`
	Bills    []string     `json:"Bills"`
	ID       string       `json:"ID"`
	Pib      string       `json:"Pib"`
	Products []string     `json:"Products"`
	Type     MerchantType `json:"Type"`
}

type MerchantType int

const (
	Supermarket MerchantType = iota
	Wholesale
	Service
)

var merchantTypeToString = map[MerchantType]string{
	Supermarket: "Supermarket",
	Wholesale:   "Wholesale",
	Service:     "Service",
}

var stringToMerchantType = map[string]MerchantType{
	"Supermarket": Supermarket,
	"Wholesale":   Wholesale,
	"Service":     Service,
}

type Product struct {
	Amount     int    `json:"Ammount"`
	ExpiryDate string `json:"ExpiryDate"`
	ID         string `json:"ID"`
	Name       string `json:"Name"`
	Price      int    `json:"Price"`
}

type User struct {
	Balance int      `json:"Balance"`
	Bills   []string `json:"Bills"`
	Email   string   `json:"Email"`
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Surname string   `json:"Surname"`
}

type Bill struct {
	Date     string `json:"Date"`
	ID       string `json:"ID"`
	Merchant string `json:"Merchant"`
	Product  string `json:"Product"`
	User     string `json:"User"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, date string) error {
	// assets := []Asset{
	// 	{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
	// 	{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
	// 	{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
	// 	{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
	// 	{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
	// 	{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	// }

	users := []User{
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

	products := []Product{
		{
			ID:   "11",
			Name: "Cheese",
			ExpiryDate: func() string {
				t, err := time.Parse(time.UnixDate, date)
				if err != nil {
					fmt.Println("Error parsing time")
					return ""
				}
				t = t.AddDate(0, 0, 5)
				return t.String()
			}(),
			Price:  5,
			Amount: 3,
		},
		{
			ID:   "12",
			Name: "Milk",
			ExpiryDate: func() string {
				t, err := time.Parse(time.UnixDate, date)
				if err != nil {
					fmt.Println("Error parsing time")
					return ""
				}
				t = t.AddDate(0, 0, 5)
				return t.String()
			}(),
			Price:  3,
			Amount: 10,
		},
		{
			ID:         "13",
			Name:       "Toilet paper",
			ExpiryDate: "",
			Price:      6,
			Amount:     20,
		},
		{
			ID:         "14",
			Name:       "Socks",
			ExpiryDate: "",
			Price:      2,
			Amount:     30,
		},
	}
	merchants := []Merchant{
		{
			ID:   "21",
			Pib:  "1234",
			Type: Wholesale,
			Products: []string{
				products[0].ID,
				products[1].ID,
			},
			Bills:   []string{},
			Balance: 1000,
		},
		{
			ID:   "22",
			Pib:  "5678",
			Type: Supermarket,
			Products: []string{
				products[2].ID,
				products[3].ID,
			},
			Bills:   []string{},
			Balance: 1000,
		},
	}

	// for _, asset := range assets {
	// 	assetJSON, err := json.Marshal(asset)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = ctx.GetStub().PutState(asset.ID, assetJSON)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to put to world state. %v", err)
	// 	}
	// }

	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(product.ID, productJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	for _, merchant := range merchants {
		merchantJSON, err := json.Marshal(merchant)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(merchant.ID, merchantJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(user.ID, userJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateMerchant(ctx contractapi.TransactionContextInterface, id string, pib string, mtype string, balance int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the merchant %s already exists", id)
	}

	merchant := Merchant{
		ID:       id,
		Pib:      pib,
		Type:     stringToMerchantType[mtype],
		Balance:  balance,
		Products: []string{},
		Bills:    []string{},
	}

	merchantJSON, err := json.Marshal(merchant)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(merchant.ID, merchantJSON)
}

func (s *SmartContract) AddProductToMerchant(ctx contractapi.TransactionContextInterface, merchantId string, productId string) error {
	merchantJSON, err := ctx.GetStub().GetState(merchantId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if merchantJSON == nil {
		return fmt.Errorf("the merchant %s does not exist", merchantId)
	}

	productJSON, err := ctx.GetStub().GetState(productId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if productJSON == nil {
		return fmt.Errorf("the product %s does not exist", productId)
	}

	var merchant Merchant
	err = json.Unmarshal(merchantJSON, &merchant)
	if err != nil {
		return err
	}

	for _, product := range merchant.Products {
		if product == productId {
			return fmt.Errorf("merchant %s already has product %s", merchantId, productId)
		}
	}

	merchant.Products = append(merchant.Products, productId)

	merchantJSON, err = json.Marshal(merchant)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(merchant.ID, merchantJSON)
}

func (s *SmartContract) CreateUsers(ctx contractapi.TransactionContextInterface, users []User) error {
	for _, user := range users {
		exists, err := s.AssetExists(ctx, user.ID)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("the user %s already exists", user.ID)
		}

		user.Bills = []string{}
		userJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(user.ID, userJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SmartContract) BuyProduct(ctx contractapi.TransactionContextInterface, userId string, productId string, date string) error {
	userJSON, err := ctx.GetStub().GetState(userId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if userJSON == nil {
		return fmt.Errorf("the user %s does not exist", userId)
	}

	productJSON, err := ctx.GetStub().GetState(productId)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if productJSON == nil {
		return fmt.Errorf("the product %s does not exist", productId)
	}

	queryString := fmt.Sprintf(`{
		"selector": {
			"Products": {
				"$in": ["%s"]
			}
		}
	}`, productId)

	merchantIter, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	defer merchantIter.Close()

	var merchantJSON *queryresult.KV
	if merchantIter.HasNext() {
		merchantJSON, err = merchantIter.Next()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("the merchant with product %s does not exist", productId)
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return err
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return err
	}

	var merchant Merchant
	err = json.Unmarshal(merchantJSON.Value, &merchant)
	if err != nil {
		return err
	}

	if user.Balance < product.Price {
		return fmt.Errorf("user does not have enough assets to buy the product")
	}

	user.Balance -= product.Price
	merchant.Balance += product.Price
	product.Amount -= 1

	currentTime, err := time.Parse(time.UnixDate, date)
	if err != nil {
		return err
	}

	bill := Bill{
		Date:     currentTime.String(),
		ID:       fmt.Sprintf("%s%s%d", userId, productId, currentTime.Unix()),
		Merchant: merchant.ID,
		Product:  productId,
		User:     userId,
	}

	user.Bills = append(user.Bills, bill.ID)
	merchant.Bills = append(merchant.Bills, bill.ID)

	if product.Amount == 0 {
		ctx.GetStub().DelState(productId)
	} else {
		productJSON, err = json.Marshal(product)
		if err != nil {
			return err
		}

		ctx.GetStub().PutState(productId, productJSON)
	}

	userJSON, err = json.Marshal(user)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(userId, userJSON)

	merchantJSON.Value, err = json.Marshal(merchant)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(merchant.ID, merchantJSON.Value)

	billJSON, err := json.Marshal(bill)
	if err != nil {
		return err
	}
	ctx.GetStub().PutState(bill.ID, billJSON)

	return nil
}

func (s *SmartContract) AddFunds(ctx contractapi.TransactionContextInterface, id string, amount int) error {
	entityJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	if entityJSON == nil {
		return fmt.Errorf("no entity with id %s exists on the legder", id)
	}

	var user User
	var merchant Merchant
	err = json.Unmarshal(entityJSON, &user)
	if err != nil {
		err = json.Unmarshal(entityJSON, &merchant)
		if err != nil {
			return fmt.Errorf("entity with id %s is neither merchant nor user", id)
		}

		merchant.Balance += amount

		merchantJSON, err := json.Marshal(merchant)
		if err != nil {
			return err
		}

		return ctx.GetStub().PutState(id, merchantJSON)
	}

	user.Balance += amount

	userJSON, err := json.Marshal(user)

	return ctx.GetStub().PutState(id, userJSON)
}

func (s *SmartContract) FindProduct(ctx contractapi.TransactionContextInterface, id, name, mtype string, price int) ([]*Product, error) {
	selector := map[string]interface{}{}

	if id != "" {
		selector["ID"] = id
	}

	if name != "" {
		selector["Name"] = name
	}

	if price != -1 {
		selector["Price"] = price
	}

	query := map[string]interface{}{
		"selector": selector,
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("query failed %v", err)
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryRespone, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		if err := json.Unmarshal(queryRespone.Value, &product); err != nil {
			return nil, err
		}

		if mtype != "" {
			merchant, err := s.FindMerchantByProductIdAndMerchantType(ctx, product.ID, mtype)
			if err != nil {
				continue
			}

			if merchant == nil {
				continue
			}

		}

		products = append(products, &product)

	}

	return products, nil
}

func (s *SmartContract) FindMerchantByProductIdAndMerchantType(ctx contractapi.TransactionContextInterface, productId string, mtype string) (*Merchant, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"Products": map[string]interface{}{
				"$in": []string{productId},
			},
			"Type": stringToMerchantType[mtype],
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	resultsIterator, err := ctx.GetStub().GetQueryResult(string(queryBytes))
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var merchant Merchant
		err = json.Unmarshal(queryResponse.Value, &merchant)

		return &merchant, nil
	}

	return nil, nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("0", "1")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User

		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)

	}

	return users, nil
}

func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("1", "2")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product

		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)

	}

	return products, nil
}

func (s *SmartContract) GetAllMerchants(ctx contractapi.TransactionContextInterface) ([]*Merchant, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("2", "3")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var merchants []*Merchant
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var merchant Merchant

		err = json.Unmarshal(queryResponse.Value, &merchant)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, &merchant)

	}

	return merchants, nil
}
