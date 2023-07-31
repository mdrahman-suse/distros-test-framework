package createcluster

import (
	"flag"
	"os"
	"testing"

	"github.com/rancher/distros-test-framework/cmd"
	"github.com/rancher/distros-test-framework/lib/cluster"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	flag.Var(&cmd.ServiceFlag.ClusterConfig.Product, "product", "Distro to create cluster and run the tests")
	flag.Var(&cmd.ServiceFlag.ClusterConfig.Destroy, "destroy", "Destroy cluster after test")
	flag.Var(&cmd.ServiceFlag.ClusterConfig.Arch, "arch", "Architecture type")
	flag.Var(&cmd.ServiceFlag.ClusterConfig.KubeConfigFile, "kubeconfig_file", "Kubeconfig file of an existing cluster")

	flag.Parse()
	os.Exit(m.Run())
}

func TestClusterCreateSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Cluster Test Suite")
}

var _ = AfterSuite(func() {
	g := GinkgoT()
	if cmd.ServiceFlag.ClusterConfig.Destroy {
		status, err := activity.DestroyCluster(g, cmd.ServiceFlag.ClusterConfig.Product.String())
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal("cluster destroyed"))
	}
})
