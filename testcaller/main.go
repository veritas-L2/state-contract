package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type CallerSmartContract struct {
	contractapi.Contract
}

func (s *CallerSmartContract) InitCaller(ctx contractapi.TransactionContextInterface) (string, error) {
	return "Initialised again", nil
}


func (s *CallerSmartContract) CallFabcar(ctx contractapi.TransactionContextInterface) (string, error) {
	params := []string{"QueryAllCars"}
	
	queryArgs := make([][]byte, len(params)) 
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	res := ctx.GetStub().InvokeChaincode("fabcar", queryArgs, "ch1")

	if res.Status != 200 {
		return "", fmt.Errorf("failed to query chaincode. got error[%d]: %s", res.Status, res.Payload)
	}

	return string(res.Payload), nil	
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(CallerSmartContract))

	if err != nil {
		fmt.Printf("Error create fabcarcaller chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcarcaller chaincode: %s", err.Error())
	}
}