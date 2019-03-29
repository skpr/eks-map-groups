package main

import (
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	cliFile      = kingpin.Flag("file", "Path to the config file").Envar("EKS_MAP_GROUPS_FILE").Required().String()
	cliNamespace = kingpin.Flag("namespace", "Namespace where the ConfigMap resides").Envar("EKS_MAP_GROUPS_NAMESPACE").Required().String()
	cliConfigMap = kingpin.Flag("configmap", "Name of the ConfigMap to update").Envar("EKS_MAP_GROUPS_CONFIGMAP").Required().String()
	cliFrequency = kingpin.Flag("frequency", "Frequency to check for updates").Envar("EKS_MAP_GROUPS_FREQUENCY").Default("15s").Duration()
)

func main() {
	kingpin.Parse()

	throttle := time.Tick(*cliFrequency)

	for {
		<-throttle

		err := sync(*cliFile, *cliNamespace, *cliConfigMap)
		if err != nil {
			panic(err)
		}
	}
}

func sync(file, namespace, name string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	var groups []MapGroup

	err = yaml.Unmarshal(f, &groups)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	awsclient := iam.New(session.New(&aws.Config{}))

	users, err := getUsers(awsclient, groups)
	if err != nil {
		return errors.Wrap(err, "failed to get user list")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get Kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to get Kubernetes client")
	}

	configmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to get ConfigMap")
	}

	data, err := yaml.Marshal(&users)
	if err != nil {
		return errors.Wrap(err, "failed to build user config")
	}

	configmap.Data["mapUsers"] = string(data)

	_, err = clientset.CoreV1().ConfigMaps(namespace).Update(configmap)
	if err != nil {
		return errors.Wrap(err, "failed to update ConfigMap")
	}

	return nil
}

// Helper function to get a list of users.
func getUsers(client *iam.IAM, groups []MapGroup) ([]MapUser, error) {
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
