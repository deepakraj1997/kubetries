package utils

import (
	"context"

	kuberrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func CreateVelero(client dynamic.Interface, spec *unstructured.Unstructured, namespace string) (unstructured.Unstructured, error) {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "konveyor.openshift.io",
		Version:  "v1alpha1",
		Resource: "veleros",
	})

	veleroResource, err := resourceClient.Namespace(namespace).Create(context.Background(), spec, metav1.CreateOptions{})
	// return err
	if err != nil && kuberrs.IsAlreadyExists(err) {
		veleroResource, err := resourceClient.Namespace(namespace).Get(context.Background(), spec.GetName(), metav1.GetOptions{})
		return *veleroResource, err
	}
	return *veleroResource, err
}

func DeleteVelero(client dynamic.Interface, namespace string, instanceName string) error {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "konveyor.openshift.io",
		Version:  "v1alpha1",
		Resource: "veleros",
	})

	err := resourceClient.Namespace(namespace).Delete(context.Background(), instanceName, metav1.DeleteOptions{})
	return err
}

func CreateMyVeleroInstance(client dynamic.Interface, namespace string, instanceName string) (unstructured.Unstructured, error) {
	var veleroSpec = unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "konveyor.openshift.io/v1alpha1",
			"kind":       "Velero",
			"metadata": map[string]interface{}{
				"name":      instanceName,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"olm_managed": true,
				"default_velero_plugins": []string{
					"aws",
					"csix",
					"openshift",
				},
				"backup_storage_locations": [](map[string]interface{}){
					map[string]interface{}{
						"config": map[string]interface{}{
							"profile": "default",
							"region":  "us-east-1",
						},
						"credentials_secret_ref": map[string]interface{}{
							"name":      "cloud-credentials",
							"namespace": "oadp-operator",
						},
						"object_storage": map[string]interface{}{
							"bucket": "deepakvelero",
							"prefix": "velero",
						},
						"name":     "default",
						"provider": "aws",
					},
				},
				"velero_feature_flags": "EnableCSI",
				"enable_restic":        true,
				"volume_snapshot_locations": [](map[string]interface{}){
					map[string]interface{}{
						"config": map[string]interface{}{
							"profile": "default",
							"region":  "us-west-1",
						},
						"name":     "default",
						"provider": "aws",
					},
				},
			},
		},
	}
	return CreateVelero(client, &veleroSpec, namespace)
}

func CreateVeleroBackup(client dynamic.Interface, spec *unstructured.Unstructured, namespace string) (unstructured.Unstructured, error) {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "velero.io",
		Version:  "v1",
		Resource: "backups",
	})

	veleroResource, err := resourceClient.Namespace(namespace).Create(context.Background(), spec, metav1.CreateOptions{})
	// return err
	if err != nil && kuberrs.IsAlreadyExists(err) {
		veleroResource, err := resourceClient.Namespace(namespace).Get(context.Background(), spec.GetName(), metav1.GetOptions{})
		return *veleroResource, err
	}
	return *veleroResource, err
}

func GetVeleroBackup(client dynamic.Interface, name string, namespace string) (unstructured.Unstructured, error) {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "velero.io",
		Version:  "v1",
		Resource: "restores",
	})
	veleroResource, err := resourceClient.Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	return *veleroResource, err
}

func RestoreVeleroBackup(client dynamic.Interface, spec *unstructured.Unstructured, namespace string) (unstructured.Unstructured, error) {
	resourceClient := client.Resource(schema.GroupVersionResource{
		Group:    "velero.io",
		Version:  "v1",
		Resource: "restores",
	})

	veleroResource, err := resourceClient.Namespace(namespace).Create(context.Background(), spec, metav1.CreateOptions{})
	// return err
	if err != nil && kuberrs.IsAlreadyExists(err) {
		return GetVeleroBackup(client, spec.GetName(), namespace)
	}
	return *veleroResource, err
}
