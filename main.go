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


func (s *StateContract) InitStateContract(ctx contractapi.TransactionContextInterface) (error){
	s.state = *mpt.NewTrie() 

	return nil
}

func (s *StateContract) PutState(ctx contractapi.TransactionContextInterface, key string, value string) (string, error){
	s.state.Put([]byte(key), []byte(value))
	res, found := s.state.Get([]byte(key))

	if (!found){
		return "", fmt.Errorf("failed to find key %s in world state", key)
	}

	return string(res), nil
}

func (s *StateContract) DeleteState(ctx contractapi.TransactionContextInterface, key string) (string, error){
	s.state.Put([]byte(key), []byte(nil))
	res, found := s.state.Get([]byte(key))

	if (!found){
		return "", fmt.Errorf("failed to find key %s in world state", key)
	}

	return string(res) , nil
}

func (s *StateContract) GetState(ctx contractapi.TransactionContextInterface, key string) (string, error){
	res, found := s.state.Get([]byte(key))

	if (!found){
		return "", fmt.Errorf("failed to find key %s in world state", key)
	}

	return string(res), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(StateContract))

	if err != nil {
		fmt.Printf("Error create statecontract chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting statecontract chaincode: %s", err.Error())
	}
}