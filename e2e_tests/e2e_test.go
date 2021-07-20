package e2e_tests

import (
	"fmt"

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

	Describe("E2E", func() {
		Context("Verify OADP Installation", func() {
			It("OADP Namespace exists", func() {
				err := utils.GetKNamespace(clientset, "nginx-example")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Verify Nginx Deployment", func() {
			err := nginx.DeployNginxStateless()
			Expect(err).NotTo(HaveOccurred())

			It("Nginx Namespace exists", func() {
				err = utils.GetKNamespace(clientset, "nginx-example")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Verify backup exists", func() {
			It("Backup Resource exists", func() {
				_, err := utils.GetVeleroBackup(dclient, "nginx-stateless", "nginx-example")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		var _ = Describe("Delete Namespace", func() {
			err := utils.DeleteKNamespace(clientset, "nginx-example")
			Expect(err).NotTo(HaveOccurred())
			// verify namespace
		})

		Context("Restore Nginx", func() {
			It("Restore Resource exists", func() {
				err := nginx.RestoreNginxStateless()
				Expect(err).NotTo(HaveOccurred())
				// verify pods
			})
		})
	})
})

var _ = Describe("Test", func() {
	fmt.Printf("hi")
})

var _ = AfterSuite(func() {
	err := DestroyTestSuite("oadp-operator")
	Expect(err).NotTo(HaveOccurred())
})
