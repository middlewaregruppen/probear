package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ProbearPods struct {
	Name string
	Addr string
	Node string
}

func GetProbearPods() ([]ProbearPods, error) {
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

	var res = make([]ProbearPods, pods.Size())

	for k, p := range pods.Items {
		res[k].Name = p.GetName()
		res[k].Addr = p.Status.PodIP
		res[k].Node = p.Spec.NodeName

	}
	fmt.Printf("Targets: %+v", res)

	return res, nil
}
