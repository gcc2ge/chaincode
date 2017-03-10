package main

import(
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
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

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response{
	function,args:=stub.GetFunctionAndParameters()

	if function != "invoke" {
                return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0]  == "write"{
		return t.write(stub,args)
	} 
	if args[0] == "queryFund"{
		return t.queryFund(stub,args)
	}

	return shim.Error("Unknown action, check the first argument, must be one of 'write', 'queryFund'")
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	id:=args[1]
	project_id :=args[2]
	fundraiser_id :=args[3]
	use_pople :=args[4]
	use_type :=args[5]
	use_nums :=args[6]
	use_dt :=args[7]
	use_desc :=args[8]
	bills :=args[9]
	bills_abstract :=args[10]
	createdt :=args[11]
	modifydt :=args[12]


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

	crowdFundBytes,_:=json.Marshal(crowdFundJson)
	err:=stub.PutState(id,crowdFundBytes)
	if err!=nil{
		return shim.Error(err.Error())
	}

	fmt.Println("write %d success",id)

	return shim.Success(nil)

}

func (t *SimpleChaincode) queryFund(stub shim.ChaincodeStubInterface,args []string) pb.Response{
	if len(args) <1{
		return shim.Error("incorrect number of arguments. Expectiong 1")
	}

	var crowdFundJson crowdFund
	id:=args[1]
	crowdFundBytes,err:=stub.GetState(id)

	if err!=nil{
		return shim.Error("failed to get marble: "+err.Error())
	}

	err=json.Unmarshal([]byte(crowdFundBytes),&crowdFundJson)
	if err!=nil{
		return shim.Error("failed to delete state: "+err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("project_id:"+crowdFundJson.Project_id)
	buffer.WriteString("fundraiser_id:"+crowdFundJson.Fundraiser_id)
	buffer.WriteString("use_pople:"+crowdFundJson.Use_pople)
	buffer.WriteString("use_type:"+crowdFundJson.Use_type)
	buffer.WriteString("use_nums:"+crowdFundJson.Use_nums)
	buffer.WriteString("use_dt:"+crowdFundJson.Use_dt)
	buffer.WriteString("use_desc:"+crowdFundJson.Use_desc)
	buffer.WriteString("bills_abstract:"+crowdFundJson.Bills_abstract)
	buffer.WriteString("createdt:"+crowdFundJson.Createdt)
	buffer.WriteString("modifydt:"+crowdFundJson.Modifydt)
	buffer.WriteString("}")
	return shim.Success(buffer.Bytes())
}

