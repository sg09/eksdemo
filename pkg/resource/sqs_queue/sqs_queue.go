package sqs_queue

import (
	"eksdemo/pkg/cmd"
	"eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "sqs-queue",
			Description: "SQS Queue",
			Aliases:     []string{"sqs-queues", "queues", "queue", "sqs"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
