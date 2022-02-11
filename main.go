package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/veritas-L2/merkle-patricia-trie/src/mpt"
)


type SCState = map[string] string

type StateContract struct {
	contractapi.Contract
	state mpt.Trie
}


func (s *StateContract) InitStateContract(ctx contractapi.TransactionContextInterface) (mpt.Trie, error){
	s.state = *mpt.NewTrie() 
	s.state.Put([]byte("Hello"), []byte("World"))

	return s.state, nil
}

func (s *StateContract) PutState(ctx contractapi.TransactionContextInterface, key string, value string) (mpt.Trie, error){
	s.state.Put([]byte(key), []byte(value))

	return s.state, nil
}

func (s *StateContract) DeleteState(ctx contractapi.TransactionContextInterface, key string) ([]byte, bool){
	s.state.Put([]byte(key), []byte(nil))

	return s.state.Get([]byte(key))
}

func (s *StateContract) GetState(ctx contractapi.TransactionContextInterface, key string) ([]byte, bool){
	return s.state.Get([]byte(key))
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