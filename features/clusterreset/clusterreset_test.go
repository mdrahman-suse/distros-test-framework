package clusterreset

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test:", func() {
	It("Start Up with no issues", func() {
		specs.TestBuildCluster(cluster)
	})

	It("Validate Nodes Before Reset", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Validate Pods Before Reset", func() {
		specs.TestPodStatus(
			cluster,
			assert.PodAssertRestart(),
			assert.PodAssertReady())
	})

	It("Verifies ClusterIP Service Before Reset", func() {
		specs.TestServiceClusterIP(true, true)
	})

	It("Verifies NodePort Service Before Reset", func() {
		specs.TestServiceNodePort(true, false)
	})

	It("Verifies Cluster Reset", func() {
		specs.TestClusterReset(cluster, k8sClient)
	})

	It("Validate Nodes After Reset", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Validate Pods After Reset", func() {
		specs.TestPodStatus(
			cluster,
			assert.PodAssertRestart(),
			assert.PodAssertReady())
	})

	It("Verifies Ingress After Reset", func() {
		specs.TestIngress(true, true)
	})

	It("Verifies Daemonset After Reset", func() {
		specs.TestDaemonset(true, true)
	})

	It("Verifies NodePort Service After Reset", func() {
		specs.TestServiceNodePort(false, true)
	})

	It("Verifies dns access After Reset", func() {
		specs.TestDNSAccess(true, true)
	})

	if cluster.Config.Product == "k3s" {
		It("Verifies Local Path Provisioner storage After Reset", func() {
			specs.TestLocalPathProvisionerStorage(cluster, true, true)
		})

		It("Verifies LoadBalancer Service After Reset", func() {
			specs.TestServiceLoadBalancer(true, true)
		})
	}
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
