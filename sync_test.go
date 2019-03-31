package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// IAMMock will mock the IAM client.
type IAMMock struct{}

// GetGroup returns a mock list of Users.
func (m IAMMock) GetGroup(*iam.GetGroupInput) (*iam.GetGroupOutput, error) {
	return &iam.GetGroupOutput{
		Users: []*iam.User{
			{
				Arn: aws.String("xxxxxxxxxxxxxxx"),
			},
		},
	}, nil
}

func TestSync(t *testing.T) {
	var (
		iamclient = IAMMock{}
		k8sclient = fake.NewSimpleClientset()
	)

	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "kube-system",
			Name:      "aws-auth",
		},
		Data: make(map[string]string),
	}

	_, err := k8sclient.CoreV1().ConfigMaps(configmap.ObjectMeta.Namespace).Create(configmap)
	assert.Nil(t, err)

	err = sync(iamclient, k8sclient, "testdata/groups.yml", configmap.ObjectMeta.Namespace, configmap.ObjectMeta.Name)
	assert.Nil(t, err)

	configmap, err = k8sclient.CoreV1().ConfigMaps(configmap.ObjectMeta.Namespace).Get(configmap.ObjectMeta.Name, metav1.GetOptions{})
	assert.Nil(t, err)

	assert.Contains(t, configmap.Data["mapUsers"], "userarn: xxxxxxxxxxxxxxx")
}
