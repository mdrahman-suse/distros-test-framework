package rebootinstances

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test:", func() {

	It("Start Up with no issues on rebootinstances test", func() {
		specs.TestBuildCluster(cluster)
	})

	It("Validate Nodes", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Validate Pods", func() {
		specs.TestPodStatus(
			cluster,
			assert.PodAssertRestart(),
			assert.PodAssertReady())
	})

	It("Verifies ClusterIP Service", func() {
		specs.TestServiceClusterIP(true, true)
	})

	It("Verifies NodePort Service", func() {
		specs.TestServiceNodePort(true, true)
	})

	It("Verifies Daemonset", func() {
		specs.TestDaemonset(true, true)
	})

	if cluster.Config.Product == "rke2" {
		It("Verifies Snapshot Webhook", func() {
			err := specs.TestSnapshotWebhook(true)
			Expect(err).To(HaveOccurred(), err)
		})
	}

	It("Reboot server and agent nodes", func() {
		specs.TestRebootInstances(cluster)
	})

	It("Validate Nodes after reboot", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Verifies ClusterIP Service after reboot", func() {
		specs.TestServiceClusterIP(true, true)
	})

	It("Verifies NodePort Service after reboot", func() {
		specs.TestServiceNodePort(true, true)
	})

	It("Verifies Daemonset after reboot", func() {
		specs.TestDaemonset(true, true)
	})

	It("Verifies dns access after reboot", func() {
		specs.TestDNSAccess(true, true)
	})

	if cluster.Config.Product == "k3s" {
		It("Verifies Local Path Provisioner storage", func() {
			specs.TestLocalPathProvisionerStorage(cluster, true, true)
		})

		It("Verifies LoadBalancer Service", func() {
			specs.TestServiceLoadBalancer(true, true)
		})

		It("Verifies Traefik IngressRoute using old GKV", func() {
			specs.TestIngressRoute(cluster, true, true, "traefik.containo.us/v1alpha1")
		})

		It("Verifies Traefik IngressRoute using new GKV", func() {
			specs.TestIngressRoute(cluster, true, true, "traefik.io/v1alpha1")
		})
	}
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n", CurrentSpecReport().FullText())
	}
})
