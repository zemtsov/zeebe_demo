package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/zemtsov/zeebe_demo/hlf"
	"github.com/zemtsov/zeebe_demo/zeebe"
	"log"
)

const (
	BrokerAddr = "0.0.0.0:26500"
	ProcessID  = "TestBPMN"
)

func main() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         BrokerAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	//processCreator := zeebe.NewProcessCreator(client, ProcessID)
	//processCreator.Do()

	listener := zeebe.NewListener(client)

	listener.AddHandler("checkKey", handleCheckKey)
	listener.AddHandler("addUser", handleAddUser)
	listener.AddHandler("error", handleError)
	listener.AddHandler("result", handleResult)

	listener.Listen()
}

func handleCheckKey(client worker.JobClient, job entities.Job) {

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failJob(client, job)
		return
	}

	userKey, ok := variables["key"].(string)
	if !ok {
		sendError("public key not provided", client, job)
		return
	}

	args := []string{
		toBase64(userKey),
	}

	variables["userExists"] = hlf.QueryChaincode("acl", "checkKeys", args...) == nil
	variables["message"] = "user already exists"

	sendVariables(variables, client, job)
}

func handleAddUser(client worker.JobClient, job entities.Job) {

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failJob(client, job)
		return
	}

	userKey, ok := variables["key"].(string)
	if !ok {
		sendError("public key not provided", client, job)
		return
	}

	kycHash, ok := variables["kyc"].(string)
	if !ok {
		sendError("KYC hash not provided", client, job)
		return
	}

	userID, ok := variables["userId"].(string)
	if !ok {
		sendError("user ID not provided", client, job)
		return
	}

	industrial, ok := variables["industrial"].(bool)
	if !ok {
		sendError("industrial flag not provided", client, job)
		return
	}

	args := []string{
		toBase64(userKey),
		toBase64(kycHash),
		toBase64(userID),
		toBase64(fmt.Sprintf("%t", industrial)),
	}

	if err = hlf.InvokeChaincode("acl", "addUser", args...); err != nil {
		sendError(err.Error(), client, job)
		return
	}

	sendMessage("user added successfully", client, job)
}

func handleResult(client worker.JobClient, job entities.Job) {
	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failJob(client, job)
		return
	}

	message, ok := variables["message"].(string)
	if !ok {
		message = ""
	}

	log.Printf("Message: %s", message)

	sendVariables(variables, client, job)
}

func handleError(client worker.JobClient, job entities.Job) {
	variables, err := job.GetVariablesAsMap()
	if err != nil {
		failJob(client, job)
		return
	}

	errMsg, ok := variables["error"].(string)
	if !ok {
		errMsg = "unknown error"
	}

	log.Printf("error occurred: %s", errMsg)

	sendVariables(variables, client, job)
}

func sendError(err string, client worker.JobClient, job entities.Job) {
	variables := make(map[string]interface{})
	variables["error"] = err
	sendVariables(variables, client, job)
}

func sendMessage(message string, client worker.JobClient, job entities.Job) {
	variables := make(map[string]interface{})
	variables["message"] = message
	sendVariables(variables, client, job)
}

func sendVariables(variables map[string]interface{}, client worker.JobClient, job entities.Job) {

	request, err := client.NewCompleteJobCommand().JobKey(job.GetKey()).VariablesFromMap(variables)
	if err != nil {
		failJob(client, job)
		return
	}

	_, err = request.Send(context.Background())
	if err != nil {
		failJob(client, job)
		return
	}
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}

func toBase64(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}
