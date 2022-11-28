package atmz

import (
	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"log"
)

func JobHandler(client zbc.Client, jobType string) worker.JobHandler {

	switch jobType {
	case JobCheckKey:
		return checkKeyJobHandler(client)
	case JobAddUser:
		return addUserJobHandler(client)
	case JobPublishCheckKeyResult:
		return publishCheckKeyResultHandler(client)
	case JobPublishAddUserResult:
		return publishAddUserResultHandler(client)
	case JobError:
		return errorHandler()
	case JobResult:
		return resultHandler()
	}
	return notImplementedJobHandler
}

func notImplementedJobHandler(worker.JobClient, entities.Job) {
	log.Println("this job handler is not implemented")
}
