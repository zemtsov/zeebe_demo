package atmz

import (
	"context"
	"encoding/json"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/zemtsov/zeebe_demo/hlf"
)

func checkKeyJobHandler(zbcClient zbc.Client) worker.JobHandler {

	throwErrorCommand := zbcClient.NewThrowErrorCommand()

	return func(client worker.JobClient, job entities.Job) {

		variables := processVariables{}
		if err := json.Unmarshal([]byte(job.GetVariables()), &variables); err != nil {
			_, _ = throwErrorCommand.JobKey(job.GetKey()).ErrorCode(ErrorCheckKey).Send(context.Background())
			failJobNoRetries(client, job)
			return
		}

		args := []string{
			toBase64(variables.Key),
		}
		variables.UserExists = hlf.QueryChaincode("acl", "checkKeys", args...) == nil

		completeJobWithVariables(variables, client, job)
	}
}

func publishCheckKeyResultHandler(zbcClient zbc.Client) worker.JobHandler {

	publishMessageCommand := zbcClient.NewPublishMessageCommand()
	throwErrorCommand := zbcClient.NewThrowErrorCommand()

	return func(client worker.JobClient, job entities.Job) {

		variables := processVariables{}
		if err := json.Unmarshal([]byte(job.GetVariables()), &variables); err != nil {
			_, _ = throwErrorCommand.JobKey(job.GetKey()).ErrorCode(ErrorCheckKey).Send(context.Background())
			failJobNoRetries(client, job)
			return
		}

		message := MsgUserNotFound
		if variables.UserExists {
			message = MsgUserFound
		}
		_, _ = publishMessageCommand.MessageName(message).CorrelationKey(variables.Key).Send(context.Background())

		completeJobWithVariables(variables, client, job)
	}
}
