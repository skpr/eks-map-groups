package main

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Helper function to sync the
func sync(iamclient IAMClient, k8sclient kubernetes.Interface, file, namespace, name string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	var groups []MapGroup

	err = yaml.Unmarshal(f, &groups)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	users, err := getUsers(iamclient, groups)
	if err != nil {
		return errors.Wrap(err, "failed to get user list")
	}

	configmap, err := k8sclient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get ConfigMap")
	}

	data, err := yaml.Marshal(&users)
	if err != nil {
		return errors.Wrap(err, "failed to build user config")
	}

	configmap.Data["mapUsers"] = string(data)

	_, err = k8sclient.CoreV1().ConfigMaps(namespace).Update(configmap)
	if err != nil {
		return errors.Wrap(err, "failed to update ConfigMap")
	}

	return nil
}

// Helper function to get a list of users.
func getUsers(client IAMClient, groups []MapGroup) ([]MapUser, error) {
	var users []MapUser

	for _, group := range groups {
		resp, err := client.GetGroup(&iam.GetGroupInput{
			GroupName: aws.String(group.Name),
		})
		if err != nil {
			return users, errors.Wrap(err, "failed to load group")
		}

		for _, user := range resp.Users {
			users = append(users, MapUser{
				UserARN:  *user.Arn,
				Username: group.Username,
				Groups:   group.Groups,
			})
		}
	}

	return users, nil
}
