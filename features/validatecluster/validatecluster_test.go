package validatecluster

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test:", func() {
	It("Start Up with no issues", func() {
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

	It("Verifies Ingress", func() {
		specs.TestIngress(true, true)
	})

	It("Verifies Daemonset", func() {
		specs.TestDaemonset(true, true)
	})

	It("Verifies dns access", func() {
		specs.TestDNSAccess(true, true)
	})

	if cluster.Config.Product == "rke2" {
		It("Verifies Snapshot Webhook", func() {
			err := specs.TestSnapshotWebhook(true)
			Expect(err).To(HaveOccurred(), err)
		})
	}

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
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
