package nginx

import (
	"fmt"

	"github.com/deepakraj1997/kubetries/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func BackupNginxStateless(createVelero bool) {
	fmt.Print("\n Backup Nginx using Go Client")
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
	if createVelero {
		_, err = utils.CreateMyVeleroInstance(dclient, namespace)
		if err != nil {
			panic(err.Error())
		}
	}
	var backupSpec = unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "velero.io/v1",
			"kind":       "Backup",
			"metadata": map[string]interface{}{
				"name":      "nginx-stateless",
				"namespace": namespace,
				// "labels":    veleroLabel,
			},
			"spec": map[string]interface{}{
				"hooks": map[string]interface{}{},
				"includedNamespaces": []string{
					"nginx-example",
				},
				"storageLocation": "default",
				"ttl":             "720h0m0s",
			},
		},
	}
	_, err = utils.CreateVeleroBackup(dclient, &backupSpec, namespace)
	if err != nil {
		panic(err.Error())
	}
}
