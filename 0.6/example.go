/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// This chaincode implements a simple map that is stored in the state.
// The following operations are available.

// Invoke operations
// put - requires two arguments, a key and value
// remove - requires a key

// Query operations
// get - requires one argument, a key, and returns a value
// keys - requires no arguments, returns all keys

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// Init is a no-op
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Invoke has two functions
// put - takes two arguements, a key and value, and stores them in the state
// remove - takes one argument, a key, and removes if from the state
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {
	case "write":
		id:=args[0]
		project_id :=args[1]
		fundraiser_id :=args[2]
		use_pople :=args[3]
		use_type :=args[4]
		use_nums :=args[5]
		use_dt :=args[6]
		use_desc :=args[7]
		bills :=args[8]
		bills_abstract :=args[9]
		createdt :=args[10]
		modifydt :=args[11]

		var json=""
		json+="{"
		json+="project_id:"+project_id+","
		json+="fundraiser_id:"+fundraiser_id+","
		json+="use_pople:"+use_pople+","
		json+="use_type:"+use_type+","
		json+="use_nums:"+use_nums+","
		json+="use_dt:"+use_dt+","
		json+="use_desc:"+use_desc+","
		json+="bills:"+bills+","
		json+="bills_abstract:"+bills_abstract+","
		json+="createdt:"+createdt+","
		json+="modifydt:"+modifydt
		json+="}"
		// return buffer.Bytes(),nil

		_, err := stub.GetState(id)
		if err!=nil{
			return nil,fmt.Errorf("no fund found %s",err)
		}

		err = stub.DelState(id) 
		if err!=nil{
			return nil,fmt.Errorf("cann't delete fund %s",err)
		}

		err=stub.PutState(id,[]byte(json))
		if err!=nil{
			return nil,fmt.Errorf("put fund error %s",err)
		}

		return nil,nil


	default:
		return nil, errors.New("Unsupported operation")
	}
}

// Query has two functions
// get - takes one argument, a key, and returns the value for the key
// keys - returns all keys stored in this chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "queryFund":
		if len(args) < 1 {
			return nil, errors.New("get operation must include one argument, a key")
		}
		key := args[0]
		value, err := stub.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
		}
		return value, nil


	default:
		return nil, errors.New("Unsupported operation")
	}
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
