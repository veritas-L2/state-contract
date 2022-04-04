package main

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Database struct {
	ctx contractapi.TransactionContextInterface
}

func NewDatabase(ctx contractapi.TransactionContextInterface) *Database {
	return &Database{ctx}
}

func (db *Database) Get(key []byte) (value []byte, err error){
	base64Key := base64.StdEncoding.EncodeToString(key) 
	val, err := db.ctx.GetStub().GetState(base64Key)

	res, _ := db.ctx.GetStub().GetState("team")

	println(res, string(res))

	if (val == nil){
		err = fmt.Errorf("key error: %s", string(key))
	}

	return val, err
}

func (db *Database) Put(key []byte, value []byte) (error) {
	base64Key := base64.StdEncoding.EncodeToString(key)

	db.ctx.GetStub().PutState("team", []byte("rocket"))
	res, err := db.ctx.GetStub().GetState("team")

	if(err != nil){
		print(fmt.Sprintf("failed to retrieve client's identity. %s \n", err.Error()))
	}

	print("hello\n", string(res))

	return db.ctx.GetStub().PutState(base64Key, value)
}

func (db *Database) Delete(key []byte) error {
	base64Key := base64.StdEncoding.EncodeToString(key) 

	return db.ctx.GetStub().DelState(base64Key)
}