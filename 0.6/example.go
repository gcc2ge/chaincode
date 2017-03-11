package main

import(
	// "bytes"
	"encoding/json"
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct{

}

type crowdFund struct{
	Project_id string `json:project_id`
	Fundraiser_id string `json:fundrasier_id`
	Use_pople string `json:use_people`
	Use_type string `json:`
	Use_nums string `json:` 
	Use_dt string `json:`
	Use_desc string `json:use_desc`
	Bills string `json:bills`
	Bills_abstract string `json:bills_abstract`
	Createdt string `json:createdt`
	Modifydt string `json:modifydt`
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


	crowdFundJson:=&crowdFund{
		Project_id:project_id,
		Fundraiser_id:fundraiser_id,
		Use_pople:use_pople,
		Use_type:use_type,
		Use_nums:use_nums,
		Use_dt:use_dt,
		Use_desc:use_desc,
		Bills:bills,
		Bills_abstract:bills_abstract,
		Createdt:createdt,
		Modifydt:modifydt,
	}

	_, err := stub.GetState(id)
	if err!=nil{
		return nil,fmt.Errorf("no fund found %s",err)
	}

	err = stub.DelState(id) 
	if err!=nil{
		return nil,fmt.Errorf("cann't delete fund %s",err)
	}

	crowdFundBytes,_:=json.Marshal(crowdFundJson)
	err=stub.PutState(id,crowdFundBytes)
	if err!=nil{
		return nil,fmt.Errorf("put fund error %s",err)
	}

	return nil,nil

}

func (t *SimpleChaincode) queryFund(stub shim.ChaincodeStubInterface,args []string) ([]byte, error){
	if len(args) <1{
		return nil, errors.New("get operation must include one argument, a key")
	}

	var crowdFundJson crowdFund
	id:=args[0]
	crowdFundBytes,err:=stub.GetState(id)

	if err!=nil{
		return nil, fmt.Errorf("failed to get fund : %s", err)
	}


	_=crowdFundBytes
	_=crowdFundJson
	return []byte("hello world"),nil

	// err=json.Unmarshal([]byte(crowdFundBytes),&crowdFundJson)
	// if err!=nil{
	// 	return nil,fmt.Errorf("unmarshal error : %s ",err)
	// }

	// var buffer bytes.Buffer
	// buffer.WriteString("{")
	// buffer.WriteString("project_id:"+crowdFundJson.Project_id+",")
	// buffer.WriteString("fundraiser_id:"+crowdFundJson.Fundraiser_id+",")
	// buffer.WriteString("use_pople:"+crowdFundJson.Use_pople+",")
	// buffer.WriteString("use_type:"+crowdFundJson.Use_type+",")
	// buffer.WriteString("use_nums:"+crowdFundJson.Use_nums+",")
	// buffer.WriteString("use_dt:"+crowdFundJson.Use_dt+",")
	// buffer.WriteString("use_desc:"+crowdFundJson.Use_desc+",")
	// buffer.WriteString("bills_abstract:"+crowdFundJson.Bills_abstract+",")
	// buffer.WriteString("createdt:"+crowdFundJson.Createdt+",")
	// buffer.WriteString("modifydt:"+crowdFundJson.Modifydt)
	// buffer.WriteString("}")
	// return buffer.Bytes(),nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "queryFund"{
		return t.queryFund(stub,args)
	}
	return nil, errors.New("Unsupported operation")		
}

