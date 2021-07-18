package utils

import (
	"context"
	"flag"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	Check(err)
	return string(dat)
}

func LoadKConfig(path string) (*kubernetes.Clientset, error) {
	var kubeconfig *string

	kubeconfig = flag.String("kubeconfig", path, "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

func LoadDConfig(path string) (dynamic.Interface, error) {
	var kubeconfig *string

	kubeconfig = flag.String("kubeconfig", path, "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := dynamic.NewForConfig(config)
	return clientset, err
}

func CreateKNamespace(clientset *kubernetes.Clientset, namespace string, labels map[string]string) error {
	ns := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: labels,
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	return err
}

func DeleteKNamespace(clientset *kubernetes.Clientset, namespace string) error {
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	return err
}

func CreateKDeployment(clientset *kubernetes.Clientset, deploymentMeta *metav1.ObjectMeta, deploymentSpec *appsv1.DeploymentSpec) error {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: *deploymentMeta,
		Spec:       *deploymentSpec,
	}
	_, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	return err
}

func CreateKService(clientset *kubernetes.Clientset, serviceMeta *metav1.ObjectMeta, serviceSpec *apiv1.ServiceSpec) error {
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	service := &apiv1.Service{
		ObjectMeta: *serviceMeta,
		Spec:       *serviceSpec,
	}
	_, err := serviceClient.Create(context.Background(), service, metav1.CreateOptions{})
	return err
}

func CreateOroute(client dynamic.Interface, spec *unstructured.Unstructured, namespace string) error {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "route.openshift.io",
		Version:  "v1",
		Resource: "routes",
	})
	_, err := resourceClient.Namespace(namespace).Create(context.Background(), spec, metav1.CreateOptions{})
	return err
}
