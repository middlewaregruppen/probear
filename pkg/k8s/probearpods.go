package k8s

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

/* ProbearPod is an represtaion of a pod. 
*  It contains the basics, name, address and location of the pod.
*/
type ProbearPod struct {
	Name   string
	Addr   string
	Node   string
	Region string
	Zone   string
}

/* GetProbearPods scans the cluster for instances of Probear in the namespace
*  that the pod is currently running in.
*  The label app=probear must be present on the pod.
*/

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

	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		namespace = []byte("probear")
	}

	// get probear pods in the current namespace
	pods, err := clientset.CoreV1().Pods(string(namespace)).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=probear"})
	if err != nil {
		return nil, err
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	var res = make([]ProbearPod, len(pods.Items))

	for k, p := range pods.Items {
		res[k].Name = p.GetName()
		res[k].Addr = p.Status.PodIP
		res[k].Node = p.Spec.NodeName
		res[k].Region = getRegion(p.Spec.NodeName)
		res[k].Zone = getZone(p.Spec.NodeName)
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
	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		namespace = []byte("probear")
	}

	// get probear pods in the probear namespace
	pod, err := clientset.CoreV1().Pods(string(namespace)).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if pod == nil {
		return nil, errors.New("pod not found")
	}

	return &ProbearPod{
		Name:   pod.GetName(),
		Node:   pod.Spec.NodeName,
		Addr:   pod.Status.PodIP,
		Region: getRegion(pod.Spec.NodeName),
		Zone:   getZone(pod.Spec.NodeName),
	}, nil

}

func getRegion(nodeName string) (region string) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return ""
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ""
	}

	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return ""
	}

	if len(node.Labels["topology.kubernetes.io/region"]) > 1 {
		return node.Labels["topology.kubernetes.io/region"]
	}

	if len(node.Labels["failure-domain.beta.kubernetes.io/region"]) > 1 {
		return node.Labels["failure-domain.beta.kubernetes.io/region"]
	}
	return ""
}

func getZone(nodeName string) (region string) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return ""
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ""
	}

	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return ""
	}

	if len(node.Labels["topology.kubernetes.io/zone"]) > 1 {
		return node.Labels["topology.kubernetes.io/zone"]
	}

	if len(node.Labels["failure-domain.beta.kubernetes.io/zone"]) > 1 {
		return node.Labels["failure-domain.beta.kubernetes.io/zone"]
	}
	return ""
}
