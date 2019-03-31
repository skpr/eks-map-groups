package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
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

		iamclient := iam.New(session.New(&aws.Config{}))

		k8sconfig, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}

		k8sclient, err := kubernetes.NewForConfig(k8sconfig)
		if err != nil {
			panic(err)
		}

		err = sync(iamclient, k8sclient, *cliFile, *cliNamespace, *cliConfigMap)
		if err != nil {
			panic(err)
		}
	}
}
