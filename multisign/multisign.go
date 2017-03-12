
package main

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("multichain")

type CrowdFundChaincode struct {
}

func (t *CrowdFundChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err := stub.CreateTable("FundOwnership", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "fund", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Owner", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating fundsOnwership table.")
	}

	adminCert, err := stub.GetCallerMetadata()
	if err != nil {
		logger.Debug("Failed getting metadata")
		return nil, errors.New("Failed getting metadata.")
	}
	if len(adminCert) == 0 {
		logger.Debug("Invalid admin certificate. Empty.")
		return nil, errors.New("Invalid admin certificate. Empty.")
	}


	stub.PutState("admin", adminCert)


	return nil, nil
}

func (t *CrowdFundChaincode) assign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fund := args[0]
	owner, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, errors.New("Failed decodinf owner")
	}

	adminCertificate, err := stub.GetState("administrator")
	if err != nil {
		return nil, errors.New("Failed fetching admin identity")
	}

	ok, err := t.isCaller(stub, adminCertificate)
	if err != nil {
		return nil, errors.New("Failed checking admin identity")
	}
	if !ok {
		return nil, errors.New("The caller is not an administrator")
	}

	logger.Debugf("New owner of [%s] is [% x]", fund, owner)

	ok, err = stub.InsertRow("FundOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: fund}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("fund was already assigned.")
	}

	return nil, err
}

func (t *CrowdFundChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fund := args[0]
	newOwner, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, fmt.Errorf("Failed decoding owner")
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: fund}}
	columns = append(columns, col1)

	row, err := stub.GetRow("FundOwnership", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving fund [%s]: [%s]", fund, err)
	}

	prvOwner := row.Columns[1].GetBytes()
	logger.Debugf("Previous owener of [%s] is [% x]", fund, prvOwner)
	if len(prvOwner) == 0 {
		return nil, fmt.Errorf("Invalid previous owner. Nil")
	}

	ok, err := t.isCaller(stub, prvOwner)
	if err != nil {
		return nil, errors.New("Failed checking fund owner identity")
	}
	if !ok {
		return nil, errors.New("The caller is not the owner of the fund")
	}

	err = stub.DeleteRow(
		"FundOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_String_{String_: fund}}},
	)
	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}

	_, err = stub.InsertRow(
		"FundOwnership",
		shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: fund}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: newOwner}},
			},
		})
	if err != nil {
		return nil, errors.New("Failed inserting row.")
	}

	logger.Debug("New owner of [%s] is [% x]", fund, newOwner)

	return nil, nil
}

func (t *CrowdFundChaincode) isCaller(stub shim.ChaincodeStubInterface, certificate []byte) (bool, error) {

	sigma, err := stub.GetCallerMetadata()
	if err != nil {
		return false, errors.New("Failed getting metadata")
	}
	payload, err := stub.GetPayload()
	if err != nil {
		return false, errors.New("Failed getting payload")
	}
	binding, err := stub.GetBinding()
	if err != nil {
		return false, errors.New("Failed getting binding")
	}

	ok, err := stub.VerifySignature(
		certificate,
		sigma,
		append(payload, binding...),
	)
	if err != nil {
		logger.Errorf("Failed checking signature [%s]", err)
		return ok, err
	}
	if !ok {
		logger.Error("Invalid signature")
	}

	return ok, err
}

func (t *CrowdFundChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "assign" {
		return t.assign(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *CrowdFundChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting 'query' but found '" + function + "'")
	}

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of an fund to query")
	}

	fund := args[0]

	logger.Debugf("Arg [%s]", string(fund))

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: fund}}
	columns = append(columns, col1)

	row, err := stub.GetRow("FundOwnership", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retriving fund [%s]: [%s]", string(fund), err)
	}

	return row.Columns[1].GetBytes(), nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(CrowdFundChaincode))
	if err != nil {
		fmt.Printf("Error starting CrowdFundChaincode: %s", err)
	}
}
