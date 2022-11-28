package atmz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/zemtsov/zeebe_demo/hlf"
)

func addUserJobHandler(zbcClient zbc.Client) worker.JobHandler {

	throwErrorCommand := zbcClient.NewThrowErrorCommand()

	return func(client worker.JobClient, job entities.Job) {

		variables := processVariables{}
		if err := json.Unmarshal([]byte(job.GetVariables()), &variables); err != nil {
			_, _ = throwErrorCommand.JobKey(job.GetKey()).ErrorCode(ErrorAddUser).Send(context.Background())
			failJobNoRetries(client, job)
			return
		}

		args := []string{
			toBase64(variables.Key),
			toBase64(variables.KYC),
			toBase64(variables.UserID),
			toBase64(fmt.Sprintf("%t", variables.Industrial)),
		}

		if err := hlf.InvokeChaincode("acl", "addUser", args...); err != nil {
			_, _ = throwErrorCommand.JobKey(job.GetKey()).ErrorCode(ErrorAddUser).Send(context.Background())
			failJob(client, job)
			return
		}

		completeJobWithVariables(variables, client, job)
	}
}

func publishAddUserResultHandler(zbcClient zbc.Client) worker.JobHandler {

	publishMessageCommand := zbcClient.NewPublishMessageCommand()
	throwErrorCommand := zbcClient.NewThrowErrorCommand()

	return func(client worker.JobClient, job entities.Job) {

		variables := processVariables{}
		if err := json.Unmarshal([]byte(job.GetVariables()), &variables); err != nil {
			_, _ = throwErrorCommand.JobKey(job.GetKey()).ErrorCode(ErrorAddUser).Send(context.Background())
			failJobNoRetries(client, job)
			return
		}

		_, _ = publishMessageCommand.MessageName(MsgUserAdded).CorrelationKey(variables.Key).Send(context.Background())

		completeJobWithVariables(variables, client, job)
	}
}
