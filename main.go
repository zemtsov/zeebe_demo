package main

import (
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/zemtsov/zeebe_demo/atmz"
	"github.com/zemtsov/zeebe_demo/zeebe"
)

const (
	BrokerAddr = "0.0.0.0:26500"
)

func main() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         BrokerAddr,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	listener := zeebe.NewListener(client)

	listener.AddHandler(atmz.JobCheckKey, atmz.JobHandler(client, atmz.JobCheckKey))
	listener.AddHandler(atmz.JobAddUser, atmz.JobHandler(client, atmz.JobAddUser))
	listener.AddHandler(atmz.JobPublishCheckKeyResult, atmz.JobHandler(client, atmz.JobPublishCheckKeyResult))
	listener.AddHandler(atmz.JobPublishAddUserResult, atmz.JobHandler(client, atmz.JobPublishAddUserResult))
	listener.AddHandler(atmz.JobError, atmz.JobHandler(client, atmz.JobError))
	listener.AddHandler(atmz.JobResult, atmz.JobHandler(client, atmz.JobResult))

	listener.Listen()
}
