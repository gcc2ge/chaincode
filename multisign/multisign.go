
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
	logger.Debug("Init Chaincode...")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	err := stub.CreateTable("FundOwnership", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "Asset", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "Owner", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating AssetsOnwership table.")
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

	logger.Debug("The administrator is [%x]", adminCert)

	stub.PutState("admin", adminCert)

	logger.Debug("Init Chaincode...done")

	return nil, nil
}

func (t *CrowdFundChaincode) assign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Assign...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	asset := args[0]
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

	logger.Debugf("New owner of [%s] is [% x]", asset, owner)

	ok, err = stub.InsertRow("FundOwnership", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: asset}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: owner}}},
	})

	if !ok && err == nil {
		return nil, errors.New("Asset was already assigned.")
	}

	logger.Debug("Assign...done!")

	return nil, err
}

func (t *CrowdFundChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Transfer...")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	asset := args[0]
	newOwner, err := base64.StdEncoding.DecodeString(args[1])
	if err != nil {
		return nil, fmt.Errorf("Failed decoding owner")
	}

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: asset}}
	columns = append(columns, col1)

	row, err := stub.GetRow("FundOwnership", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving asset [%s]: [%s]", asset, err)
	}

	prvOwner := row.Columns[1].GetBytes()
	logger.Debugf("Previous owener of [%s] is [% x]", asset, prvOwner)
	if len(prvOwner) == 0 {
		return nil, fmt.Errorf("Invalid previous owner. Nil")
	}

	ok, err := t.isCaller(stub, prvOwner)
	if err != nil {
		return nil, errors.New("Failed checking asset owner identity")
	}
	if !ok {
		return nil, errors.New("The caller is not the owner of the asset")
	}

	err = stub.DeleteRow(
		"FundOwnership",
		[]shim.Column{shim.Column{Value: &shim.Column_String_{String_: asset}}},
	)
	if err != nil {
		return nil, errors.New("Failed deliting row.")
	}

	_, err = stub.InsertRow(
		"FundOwnership",
		shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: asset}},
				&shim.Column{Value: &shim.Column_Bytes{Bytes: newOwner}},
			},
		})
	if err != nil {
		return nil, errors.New("Failed inserting row.")
	}

	logger.Debug("New owner of [%s] is [% x]", asset, newOwner)

	logger.Debug("Transfer...done")

	return nil, nil
}

func (t *CrowdFundChaincode) isCaller(stub shim.ChaincodeStubInterface, certificate []byte) (bool, error) {
	logger.Debug("Check caller...")


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

	logger.Debugf("passed certificate [% x]", certificate)
	logger.Debugf("passed sigma [% x]", sigma)
	logger.Debugf("passed payload [% x]", payload)
	logger.Debugf("passed binding [% x]", binding)

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

	logger.Debug("Check caller...Verified!")

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
	logger.Debugf("Query [%s]", function)

	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting 'query' but found '" + function + "'")
	}

	var err error

	if len(args) != 1 {
		logger.Debug("Incorrect number of arguments. Expecting name of an asset to query")
		return nil, errors.New("Incorrect number of arguments. Expecting name of an asset to query")
	}

	asset := args[0]

	logger.Debugf("Arg [%s]", string(asset))

	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: asset}}
	columns = append(columns, col1)

	row, err := stub.GetRow("FundOwnership", columns)
	if err != nil {
		logger.Debugf("Failed retriving asset [%s]: [%s]", string(asset), err)
		return nil, fmt.Errorf("Failed retriving asset [%s]: [%s]", string(asset), err)
	}

	logger.Debugf("Query done [% x]", row.Columns[1].GetBytes())

	return row.Columns[1].GetBytes(), nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(CrowdFundChaincode))
	if err != nil {
		fmt.Printf("Error starting CrowdFundChaincode: %s", err)
	}
}
