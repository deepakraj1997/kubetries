package e2e_tests

import (
	"github.com/deepakraj1997/kubetries/nginx"
	"github.com/deepakraj1997/kubetries/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/env"
)

var dclient dynamic.Interface
var clientset *kubernetes.Clientset

var _ = BeforeSuite(func() {
	var err error
	clientset, dclient, err = utils.LoadConfig(env.GetString("KUBECTL", ""))
	Expect(err).NotTo(HaveOccurred())
	err = SetupTestSuite("oadp-operator")
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("OADP", func() {
	var _ = Describe("Verify OADP namespace is created", func() {
		err := utils.GetKNamespace(clientset, "nginx-example")
		Expect(err).NotTo(HaveOccurred())
		// verify pods running
	})

	var _ = Describe("Verify Nginx Namespace is created after deploy", func() {
		err := nginx.DeployNginxStateless()
		Expect(err).NotTo(HaveOccurred())
		err = utils.GetKNamespace(clientset, "nginx-example")
		Expect(err).NotTo(HaveOccurred())
	})

	var _ = Describe("Verify backup exists", func() {
		_, err := utils.GetVeleroBackup(dclient, "nginx-stateless", "nginx-example")
		Expect(err).NotTo(HaveOccurred())
		// verify pods
	})

	var _ = Describe("Delete Namespace", func() {
		err := utils.DeleteKNamespace(clientset, "nginx-example")
		Expect(err).NotTo(HaveOccurred())
		// verify namespace
	})

	var _ = Describe("Restore Backup & Check States", func() {
		err := nginx.RestoreNginxStateless()
		Expect(err).NotTo(HaveOccurred())
		// verify pods
	})

})

var _ = AfterSuite(func() {
	err := DestroyTestSuite("oadp-operator")
	Expect(err).NotTo(HaveOccurred())
})
