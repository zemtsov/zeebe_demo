package atmz

import (
	"encoding/json"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"log"
)

func errorHandler() worker.JobHandler {

	return func(client worker.JobClient, job entities.Job) {

		variables := processVariables{}
		if err := json.Unmarshal([]byte(job.GetVariables()), &variables); err != nil {
			failJobNoRetries(client, job)
			return
		}

		log.Printf("error occurred: %s", variables.Error)

		completeJobWithVariables(variables, client, job)
	}
}
