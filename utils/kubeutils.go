package utils

import (
	"context"
	"flag"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	kuberrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func LoadConfig(path string) (*kubernetes.Clientset, dynamic.Interface, error) {
	var kubeconfig *string

	kubeconfig = flag.String("kubeconfig", path, "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	dclientset, err := dynamic.NewForConfig(config)
	return clientset, dclientset, err
}

// func LoadDConfig(path string) (dynamic.Interface, error) {
// 	var dkubeconfig *string

// 	dkubeconfig = flag.String("kubeconfig", path, "kubeconfig file")
// 	flag.Parse()
// 	dconfig, err := clientcmd.BuildConfigFromFlags("", *dkubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	dclientset, err := dynamic.NewForConfig(dconfig)
// 	return dclientset, err
// }

func GetKNamespace(clientset *kubernetes.Clientset, namespace string) error {
	_, exists := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	// fmt.Println(ns_exists, exists_error)
	return exists
}

func CreateKNamespace(clientset *kubernetes.Clientset, namespace string, labels map[string]string) error {
	ns := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: labels,
		},
	}

	namespaceExists := GetKNamespace(clientset, namespace)
	if namespaceExists == nil {
		return namespaceExists
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	return err
}

func DeleteKNamespace(clientset *kubernetes.Clientset, namespace string) error {
	err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
	return err
}

func CreateKDeployment(clientset *kubernetes.Clientset, deploymentMeta *metav1.ObjectMeta, deploymentSpec *appsv1.DeploymentSpec) error {
	deploymentsClient := clientset.AppsV1().Deployments(deploymentMeta.Namespace)
	deployment := &appsv1.Deployment{
		ObjectMeta: *deploymentMeta,
		Spec:       *deploymentSpec,
	}
	// _, exists := deploymentsClient.Get(context.Background(), deploymentMeta.Name, metav1.GetOptions{})
	// fmt.Println(ns_exists, exists_error)
	// if exists == nil {
	// 	// fmt.Println("here")
	// err := deploymentsClient.Delete(context.Background(), deploymentMeta.Name, metav1.DeleteOptions{})
	// }
	// if <-ch != nil {
	_, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	// }
	if err != nil && kuberrs.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func CreateKService(clientset *kubernetes.Clientset, serviceMeta *metav1.ObjectMeta, serviceSpec *apiv1.ServiceSpec) error {
	serviceClient := clientset.CoreV1().Services(serviceMeta.Namespace)
	service := &apiv1.Service{
		ObjectMeta: *serviceMeta,
		Spec:       *serviceSpec,
	}
	// _, exists := serviceClient.Get(context.Background(), serviceMeta.Name, metav1.GetOptions{})
	// // fmt.Println(ns_exists, exists_error)
	// if exists == nil {
	// 	// fmt.Println("here")
	// 	_ = serviceClient.Delete(context.Background(), serviceMeta.Name, metav1.DeleteOptions{})
	// }
	_, err := serviceClient.Create(context.Background(), service, metav1.CreateOptions{})
	// return err
	if err != nil && kuberrs.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func CreateOroute(client dynamic.Interface, spec *unstructured.Unstructured, namespace string) error {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "route.openshift.io",
		Version:  "v1",
		Resource: "routes",
	})
	// _, exists := resourceClient.Namespace(namespace).Get(context.Background(), spec.GetName(), metav1.GetOptions{})
	// // fmt.Println(ns_exists, exists_error)
	// if exists == nil {
	// 	// fmt.Println("here")
	// 	_ = resourceClient.Namespace(namespace).Delete(context.Background(), spec.GetName(), metav1.DeleteOptions{})
	// }
	_, err := resourceClient.Namespace(namespace).Create(context.Background(), spec, metav1.CreateOptions{})
	// return err
	if err != nil && kuberrs.IsAlreadyExists(err) {
		return nil
	}
	return err
}
