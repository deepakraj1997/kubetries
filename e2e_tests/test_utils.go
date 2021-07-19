package e2e_tests

import (
	"github.com/deepakraj1997/kubetries/utils"
	"k8s.io/utils/env"
)

func SetupTestSuite(operator string) error {
	var namespace string = operator + "_e2e_tests" + env.GetString("BUILD_NUMBER", "100")
	var e2eLabel map[string]string

	e2eLabel = make(map[string]string)
	e2eLabel["env"] = "test"
	e2eLabel["operator"] = operator // p

	_, dclient, err := utils.LoadConfig(env.GetString("KUBECTL", ""))
	if err != nil {
		return err
	}

	// err = utils.CreateKNamespace(clientset, namespace, e2eLabel)
	// if err != nil {
	// 	panic(err.Error())
	// }

	_, err = utils.CreateMyVeleroInstance(dclient, namespace, namespace)
	if err != nil {
		return err
	}
	return nil
}

func DestroyTestSuite(operator string) error {
	var namespace string = operator + "_e2e_tests" + env.GetString("BUILD_NUMBER", "100")
	var e2eLabel map[string]string

	e2eLabel = make(map[string]string)
	e2eLabel["env"] = "test"
	e2eLabel["operator"] = operator // p

	_, dclient, err := utils.LoadConfig(env.GetString("KUBECTL", ""))
	if err != nil {
		return err
	}

	// err = utils.CreateKNamespace(clientset, namespace, e2eLabel)
	// if err != nil {
	// 	panic(err.Error())
	// }

	err = utils.DeleteVelero(dclient, namespace, namespace)
	if err != nil {
		return err
	}
	return nil
}
