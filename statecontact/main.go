package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


type SCState = map[string] string

type StateContract struct {
	contractapi.Contract
	state SCState
}


func (s *StateContract) InitStateContract(ctx contractapi.TransactionContextInterface) (SCState, error){
	s.state = make(map[string]string)
	s.state["hello"] = "world";
	
	return s.state, nil
}

func (s *StateContract) PutState(ctx contractapi.TransactionContextInterface, key string, value string) (SCState, error){
	s.state[key] = value
	return s.state, nil
}

func (s *StateContract) DeleteState(ctx contractapi.TransactionContextInterface, key string) (SCState, error){
	delete(s.state, key)
	return s.state, nil
}

func (s *StateContract) GetState(ctx contractapi.TransactionContextInterface, key string) (string, error){
	return s.state[key], nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(StateContract))

	fmt.Println("Starting chaincode...")

	if err != nil {
		fmt.Printf("Error create statecontract chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting statecontract chaincode: %s", err.Error())
	}
}

/*
Memory remains active till contract crashes it runs an instance of the struct
*/