package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Account struct {
    DealerID    string `json:"dealerID"`
    MSISDN      string `json:"msisdn"`
    MPIN        string `json:"mpin"`
    Balance     string `json:"balance"`
    Status      string `json:"status"`
    TransAmount string `json:"transAmount"`
    TransType   string `json:"transType"`
    Remarks     string `json:"remarks"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    accounts := []Account{
        {DealerID: "D001", MSISDN: "1234567890", MPIN: "1234", Balance: "1000", Status: "active", TransAmount: "0", TransType: "credit", Remarks: "initial"},
        {DealerID: "D002", MSISDN: "0987654321", MPIN: "5678", Balance: "2000", Status: "active", TransAmount: "0", TransType: "credit", Remarks: "initial"},
    }

    for _, account := range accounts {
        accountJSON, err := json.Marshal(account)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(account.DealerID, accountJSON)
        if err != nil {
            return fmt.Errorf("failed to put account %s: %v", account.DealerID, err)
        }
    }

    return nil
}

func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, dealerID, msisdn, mpin, balance, status, transAmount, transType, remarks string) error {
    account := Account{
        DealerID:    dealerID,
        MSISDN:      msisdn,
        MPIN:        mpin,
        Balance:     balance,
        Status:      status,
        TransAmount: transAmount,
        TransType:   transType,
        Remarks:     remarks,
    }
    accountJSON, err := json.Marshal(account)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(dealerID, accountJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error create smart contract: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting smart contract: %s", err.Error())
    }
}
