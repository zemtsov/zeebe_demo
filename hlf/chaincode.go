package hlf

import (
	"encoding/json"
	"fmt"
)

type chaincodeParameters struct {
	Channel   string   `json:"channel"`
	Chaincode string   `json:"chaincodeId"`
	Function  string   `json:"fcn"`
	Arguments []string `json:"args"`
}

func InvokeChaincode(name string, fcn string, args ...string) error {
	return callChaincode(sendInvokeRequest, name, fcn, args...)
}

func QueryChaincode(name string, fcn string, args ...string) error {
	return callChaincode(sendQueryRequest, name, fcn, args...)
}

func callChaincode(call func([]byte) error, name string, fcn string, args ...string) error {

	params := &chaincodeParameters{
		Channel:   name,
		Chaincode: name,
		Function:  fcn,
		Arguments: args,
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	fmt.Println(string(paramBytes))

	return call(paramBytes)
}
