package main

import(
	"bytes"
	// "encoding/json"
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct{

}

func main(){
	err:=shim.Start(new(SimpleChaincode))
	if err!=nil{
		fmt.Printf("error starting simple chaincode: %s",err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface,function string,args []string) ([]byte,error){

	if function  == "write"{
		return t.write(stub,args)
	} 

	return nil, errors.New("Unsupported operation")
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface,args []string) ([]byte,error){
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
	json+="bills_abstract:"h+bills_abstract+","
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

}

func (t *SimpleChaincode) queryFund(stub shim.ChaincodeStubInterface,args []string) ([]byte, error){
	if len(args) <1{
		return nil, errors.New("get operation must include one argument, a key")
	}

	id:=args[0]
	crowdFundBytes,err:=stub.GetState(id)

	if err!=nil{
		return nil, fmt.Errorf("failed to get fund : %s", err)
	}


	return crowdFundBytes,nil

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "queryFund"{
		return t.queryFund(stub,args)
	}
	return nil, errors.New("Unsupported operation")		
}

