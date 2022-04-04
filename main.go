package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	mpt "github.com/veritas-L2/merkle-patricia-trie/src"
)


type SCState = map[string] string

type StateContract struct {
	contractapi.Contract
	state mpt.Trie
	lockOwner []byte
	db Database
}

func (s *StateContract) InitStateContract(ctx contractapi.TransactionContextInterface) (error){
	if(s.lockOwner != nil){
		return fmt.Errorf("failed to acquire lock on state contract")
	}
	
	client, err  := ctx.GetClientIdentity().GetID()
	if (err != nil){
		return fmt.Errorf("failed to retrieve client's identity. %s", err.Error())
	}
	s.lockOwner = []byte(client)

	s.state = *mpt.NewTrie(mpt.MODE_NORMAL)
	s.db = *NewDatabase(ctx)
	s.state.LoadFromDB(&s.db)

	return nil
}

func (s *StateContract) PutState(ctx contractapi.TransactionContextInterface, key string, value string) (string, error){
	client, err  := ctx.GetClientIdentity().GetID()
	if (err != nil){
		return "", fmt.Errorf("failed to retrieve client's identity. %s", err.Error())
	}
	
	if (s.lockOwner == nil || !bytes.Equal(s.lockOwner, []byte(client))){
		return "", fmt.Errorf("failed to put state. lock not acquired by client")
	}
	
	s.state.Put([]byte(key), []byte(value))

	res := s.state.Get([]byte(key))
	if (res == nil){
		return "", fmt.Errorf("failed to find key %s in contract state", key)
	}

	return string(res), nil
}

func (s *StateContract) DeleteState(ctx contractapi.TransactionContextInterface, key string) (string, error){
	client, err  := ctx.GetClientIdentity().GetID()
	if (err != nil){
		return "", fmt.Errorf("failed to retrieve client's identity. %s", err.Error())
	}
	
	if (s.lockOwner == nil || !bytes.Equal(s.lockOwner, []byte(client))){
		return "", fmt.Errorf("failed to delete state. lock not acquired by client")
	}
	
	res := s.state.Get([]byte(key))
	if (res == nil){
		return "", fmt.Errorf("failed to find key %s in contract state", key)
	}

	s.state.Put([]byte(key), nil)

	return string(res) , nil
}

func (s *StateContract) GetState(ctx contractapi.TransactionContextInterface, key string) (string, error){
	client, err  := ctx.GetClientIdentity().GetID()
	if (err != nil){
		return "", fmt.Errorf("failed to retrieve client's identity. %s", err.Error())
	}
	
	if (s.lockOwner == nil || !bytes.Equal(s.lockOwner, []byte(client))){
		return "", fmt.Errorf("failed to get state. lock not acquired by client")
	}
	
	res := s.state.Get([]byte(key))
	if (res == nil){
		return "", fmt.Errorf("failed to find key %s in contract state", key)
	}

	return string(res), nil
}

func (s* StateContract) ReleaseStateContract(ctx contractapi.TransactionContextInterface) (error){
	client, err  := ctx.GetClientIdentity().GetID()
	if (err != nil){
		return fmt.Errorf("failed to retrieve client's identity. %s", err.Error())
	}
	
	if (s.lockOwner == nil || !bytes.Equal(s.lockOwner, []byte(client))){
		return fmt.Errorf("failed to release state contract. lock not acquired by client")
	}

	s.state.SaveToDB(&s.db)
	
	s.lockOwner = nil
	s.state = *mpt.NewTrie(mpt.MODE_NORMAL)

	return nil
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