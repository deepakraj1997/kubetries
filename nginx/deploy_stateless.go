package nginx

import (
	"fmt"

	"github.com/deepakraj1997/kubetries/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func int32Ptr(i int32) *int32 {
	return &i
}

func DeployNginxStateless() {
	fmt.Print("\nDeploy Nginx using Go Client")
	var nginxLabel map[string]string
	nginxLabel = make(map[string]string)
	var namespace string = "nginx-example1"
	var clientset *kubernetes.Clientset
	var err error
	nginxLabel["app"] = "nginx"
	var configPath string = "/Users/drajds/.agnosticd/drajds0714ocp4b/ocp4-workshop_drajds0714ocp4b_kubeconfig"
	clientset, err = utils.LoadKConfig(configPath)
	if err != nil {
		panic(err.Error())
	}
	err = utils.CreateKNamespace(clientset, namespace, nginxLabel)
	if err != nil {
		panic(err.Error())
	}
	var deploymentMeta = metav1.ObjectMeta{
		Name:      "nginx-deployment",
		Namespace: "nginx-example",
	}
	var deploymentSpec = appsv1.DeploymentSpec{
		Replicas: int32Ptr(2),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "nginx",
			},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": "nginx",
				},
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  "nginx",
						Image: "docker.io/bitnami/nginx",
						Ports: []apiv1.ContainerPort{
							{
								Name:          "http",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 8080,
							},
						},
					},
				},
			},
		},
	}
	err = utils.CreateKDeployment(clientset, &deploymentMeta, &deploymentSpec)
	if err != nil {
		panic(err.Error())
	}
	var serviceMeta = metav1.ObjectMeta{
		Name:      "nginx-deployment",
		Namespace: "nginx-example",
		Labels: map[string]string{
			"app": "nginx",
		},
	}
	var serviceSpec = apiv1.ServiceSpec{
		Ports: []apiv1.ServicePort{
			{
				Name:     "http",
				Protocol: apiv1.ProtocolTCP,
				Port:     8080,
			},
		},
		Selector: map[string]string{
			"app": "nginx",
		},
		Type: "LoadBalancer",
	}
	err = utils.CreateKService(clientset, &serviceMeta, &serviceSpec)
	if err != nil {
		panic(err.Error())
	}
	var dclient dynamic.Interface
	dclient, err = utils.LoadDConfig(configPath)
	if err != nil {
		panic(err.Error())
	}
	var routeSpec = unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":      "nginx-deployment",
				"namespace": "nginx-example",
				"labels": map[string]string{
					"app":     "nginx",
					"service": "nginx",
				},
			},
			"spec": map[string]interface{}{
				"to": map[string]interface{}{
					"kind": "Service",
					"name": "nginx",
				},
				"port": map[string]interface{}{
					"targetPort": 8080,
				},
			},
		},
	}
	err = utils.CreateOroute(dclient, &routeSpec, namespace)
	if err != nil {
		panic(err.Error())
	}
}
