package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetProbearPods() ([]string, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// get probear pods in the probear namespace
	pods, err := clientset.CoreV1().Pods("probear").List(context.TODO(), metav1.ListOptions{LabelSelector: "app=probear"})
	if err != nil {
		return nil, err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	var res = make([]string, pods.Size())

	for k, p := range pods.Items {
		res[k] = p.GetName()
	}

	return res, nil
}
