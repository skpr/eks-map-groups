package main

import "github.com/aws/aws-sdk-go/service/iam"

type IAM interface {
	GetGroup(*iam.GetGroupInput) (*iam.GetGroupOutput, error)
}
