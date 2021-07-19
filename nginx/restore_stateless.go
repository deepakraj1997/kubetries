package nginx

import (
	"fmt"

	"github.com/deepakraj1997/kubetries/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func RestoreNginxStateless() {
	fmt.Print("\n Restore Nginx using Go Client")
	// var veleroLabel map[string]string
	// veleroLabel = make(map[string]string)
	var namespace string = "oadp-operator"
	// var clientset *kubernetes.Clientset
	var err error
	// veleroLabel["velero.io/storage-location"] = "default"
	var configPath string = "/Users/drajds/.agnosticd/drajds0714ocp4b/ocp4-workshop_drajds0714ocp4b_kubeconfig"
	var dclient dynamic.Interface
	_, dclient, err = utils.LoadConfig(configPath)
	if err != nil {
		panic(err.Error())
	}
	// _, err = utils.CreateMyVeleroInstance(dclient, namespace)
	// if err != nil {
	// 	panic(err.Error())
	// }
	var RestoreSpec = unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "velero.io/v1",
			"kind":       "Restore",
			"metadata": map[string]interface{}{
				"name":      "nginx-stateless",
				"namespace": namespace,
				// "labels":    veleroLabel,
			},
			"spec": map[string]interface{}{
				"hooks": map[string]interface{}{},
				"excludedResources": []string{
					"nodes",
					"events",
					"events.events.k8s.io",
					"backups.velero.io",
					"restores.velero.io",
					"resticrepositories.velero.io",
				},
				"backupName": "nginx-stateless",
				"restorePVs": true,
			},
		},
	}
	_, err = utils.RestoreVeleroBackup(dclient, &RestoreSpec, namespace)
	if err != nil {
		panic(err.Error())
	}
}
