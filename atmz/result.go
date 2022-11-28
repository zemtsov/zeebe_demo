package atmz

import (
	"encoding/json"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"log"
)

type resultVariables struct {
	Message string `json:"message"`
}

func resultHandler() worker.JobHandler {

	return func(client worker.JobClient, job entities.Job) {

		variables := resultVariables{}
		if err := json.Unmarshal([]byte(job.GetVariables()), &variables); err != nil {
			failJobNoRetries(client, job)
			return
		}

		log.Printf("Result: %s", variables.Message)

		completeJobWithVariables(variables, client, job)
	}
}
