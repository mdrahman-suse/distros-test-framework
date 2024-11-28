//go:build upgradereplacement

package upgradecluster

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"
	"github.com/rancher/distros-test-framework/pkg/specs/support"
)

var _ = Describe("Upgrade Node Replacement Test:", Ordered, func() {
	It("Start Up with no issues", func() {
		specs.TestBuildCluster(cluster)
	})

	It("Validate Node", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil)
	})

	It("Validate Pod", func() {
		specs.TestPodStatus(
			cluster,
			assert.PodAssertRestart(),
			assert.PodAssertReady())
	})

	It("Verifies ClusterIP Service pre-upgrade", func() {
		specs.TestServiceClusterIP(true, false)
	})

	if cluster.Config.Product == "k3s" {
		It("Verifies LoadBalancer Service pre-upgrade", func() {
			specs.TestServiceLoadBalancer(true, false)
		})
	}

	It("Verifies Ingress pre-upgrade", func() {
		specs.TestIngress(true, false)
	})

	It("Upgrade by Node replacement", func() {
		specs.TestUpgradeReplaceNode(cluster, flags)
	})

	It("Checks Node Status after upgrade and validate version", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			assert.NodeAssertVersionTypeUpgrade(flags))
	})

	It("Checks Pod Status after upgrade", func() {
		specs.TestPodStatus(
			cluster,
			assert.PodAssertRestart(),
			assert.PodAssertReady())
	})

	It("Verifies ClusterIP Service after upgrade", func() {
		specs.TestServiceClusterIP(false, true)
	})

	It("Verifies NodePort Service after upgrade applying and deleting workload", func() {
		specs.TestServiceNodePort(true, true)
	})

	It("Verifies Ingress after upgrade", func() {
		specs.TestIngress(false, true)
	})

	if cluster.Config.Product == "k3s" {
		It("Verifies LoadBalancer Service after upgrade", func() {
			specs.TestServiceLoadBalancer(false, true)
		})
	}

	AfterAll(func() {
		if flags.Destroy {
			support.DeleteEC2Nodes(cluster)
		}
	})
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
