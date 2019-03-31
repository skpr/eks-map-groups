package main

import "github.com/aws/aws-sdk-go/service/iam"

// IAMClient for interacting with AWS.
type IAMClient interface {
	GetGroup(*iam.GetGroupInput) (*iam.GetGroupOutput, error)
}
