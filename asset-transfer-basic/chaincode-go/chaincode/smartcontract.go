package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an asset
type SmartContract struct {
    contractapi.Contract
}

// Asset describes basic details of what makes up an asset
type Asset struct {
    DEALERID      string `json:"dealerId"`
    MSISDN        string `json:"msisdn"`
    MPIN          string `json:"mpin"`
    BALANCE       int    `json:"balance"`
    STATUS        string `json:"status"`
    TRANSAMOUNT   int    `json:"transAmount"`
    TRANSTYPE     string `json:"transType"`
    REMARKS       string `json:"remarks"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []Asset{
        {DEALERID: "D001", MSISDN: "1234567890", MPIN: "1111", BALANCE: 1000, STATUS: "Active", TRANSAMOUNT: 0, TRANSTYPE: "", REMARKS: ""},
        {DEALERID: "D002", MSISDN: "9876543210", MPIN: "2222", BALANCE: 500, STATUS: "Inactive", TRANSAMOUNT: 0, TRANSTYPE: "", REMARKS: ""},
    }

    for _, asset := range assets {
        assetJSON, err := json.Marshal(asset)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(asset.DEALERID, assetJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dealerId, msisdn, mpin string, balance int, status string) error {
    exists, err := s.AssetExists(ctx, dealerId)
    if err != nil {
        return err
    }
    if exists {
        return fmt.Errorf("the asset %s already exists", dealerId)
    }

    asset := Asset{
        DEALERID:      dealerId,
        MSISDN:        msisdn,
        MPIN:          mpin,
        BALANCE:       balance,
        STATUS:        status,
        TRANSAMOUNT:   0,
        TRANSTYPE:     "",
        REMARKS:       "",
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(asset.DEALERID, assetJSON)
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
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerId, msisdn, mpin string, balance int, status string) error {
    exists, err := s.AssetExists(ctx, dealerId)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the asset %s does not exist", dealerId)
    }

    // overwriting original asset with new asset
    asset := Asset{
        DEALERID:      dealerId,
        MSISDN:        msisdn,
        MPIN:          mpin,
        BALANCE:       balance,
        STATUS:        status,
        TRANSAMOUNT:   0,
        TRANSTYPE:     "",
        REMARKS:       "",
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(asset.DEALERID, assetJSON)
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
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id, newDealerId string) error {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
        return err
    }

    asset.DEALERID = newDealerId
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
    // range query with empty string for startKey and endKey does an
    // open-ended query of all assets in the chaincode namespace.
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var assets []*Asset
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset Asset
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets =append(assets, &asset)
       }

       return assets, nil
   }

   // UpdateBalance updates the balance of an asset in the world state.
   func (s *SmartContract) UpdateBalance(ctx contractapi.TransactionContextInterface, id string, newBalance int, transAmount int, transType, remarks string) error {
       asset, err := s.ReadAsset(ctx, id)
       if err != nil {
           return err
       }

       asset.BALANCE = newBalance
       asset.TRANSAMOUNT = transAmount
       asset.TRANSTYPE = transType
       asset.REMARKS = remarks
       assetJSON, err := json.Marshal(asset)
       if err != nil {
           return err
       }

       return ctx.GetStub().PutState(id, assetJSON)
   }

   // GetAssetHistory returns the history of transactions for a given asset
   func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, id string) ([]*Asset, error) {
       resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
       if err != nil {
           return nil, err
       }
       defer resultsIterator.Close()

       var assets []*Asset
       for resultsIterator.HasNext() {
           queryResponse, err := resultsIterator.Next()
           if err != nil {
               return nil, err
           }

           var asset Asset
           err = json.Unmarshal(queryResponse.Value, &asset)
           if err != nil {
               return nil, err
           }
           assets = append(assets, &asset)
       }

       return assets, nil
   }

   func main() {
       assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
       if err != nil {
           log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
       }

       if err := assetChaincode.Start(); err != nil {
           log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
       }
   }