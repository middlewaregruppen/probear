package k8s

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ProbearPod struct {
	Name string
	Addr string
	Node string
}

func GetProbearPods() ([]ProbearPod, error) {
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

	var res = make([]ProbearPod, len(pods.Items))

	for k, p := range pods.Items {
		res[k].Name = p.GetName()
		res[k].Addr = p.Status.PodIP
		res[k].Node = p.Spec.NodeName

	}
	fmt.Printf("Targets: %+v", res)

	return res, nil
}

func GetPod(name string) (*ProbearPod, error) {

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
	pod, err := clientset.CoreV1().Pods("probear").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if pod == nil {
		return nil, errors.New("pod not found")
	}

	return &ProbearPod{
		Name: pod.GetName(),
		Node: pod.Spec.NodeName,
		Addr: pod.Status.PodIP,
	}, nil

}
