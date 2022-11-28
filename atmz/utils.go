package atmz

import (
	"context"
	"encoding/base64"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"log"
)

func failJob(client worker.JobClient, job entities.Job) {
	failJobWithRetries(client, job, job.GetRetries()-1)
}

func failJobNoRetries(client worker.JobClient, job entities.Job) {
	failJobWithRetries(client, job, 0)
}

func failJobWithRetries(client worker.JobClient, job entities.Job, retries int32) {
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(retries).Send(context.Background())
	if err != nil {
		log.Printf("error occurred: %s", err.Error())
	}
}

func completeJobWithVariables(variables interface{}, client worker.JobClient, job entities.Job) {

	request, err := client.NewCompleteJobCommand().JobKey(job.GetKey()).VariablesFromObject(variables)
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

func toBase64(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}
