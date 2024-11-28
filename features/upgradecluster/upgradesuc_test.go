//go:build upgradesuc

package upgradecluster

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SUC Upgrade Tests:", func() {

	It("Starts up with no issues", func() {
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

	It("Verifies ClusterIP Service pre-upgrade", func() {
		specs.TestServiceClusterIP(true, false)
	})

	It("Verifies NodePort Service pre-upgrade", func() {
		specs.TestServiceNodePort(true, false)
	})

	It("Verifies Ingress pre-upgrade", func() {
		specs.TestIngress(true, false)
	})

	It("Verifies Daemonset pre-upgrade", func() {
		specs.TestDaemonset(true, false)
	})

	It("Verifies DNS Access pre-upgrade", func() {
		specs.TestDNSAccess(true, false)
	})

	if cluster.Config.Product == "rke2" {
		It("Verifies Snapshot Webhook pre-upgrade", func() {
			err := specs.TestSnapshotWebhook(true)
			Expect(err).To(HaveOccurred())
		})
	}

	if cluster.Config.Product == "k3s" {
		It("Verifies LoadBalancer Service before upgrade", func() {
			specs.TestServiceLoadBalancer(true, false)
		})

		It("Verifies Local Path Provisioner storage before upgrade", func() {
			specs.TestLocalPathProvisionerStorage(cluster, true, false)
		})

		It("Verifies Traefik IngressRoute before upgrade using old GKV", func() {
			specs.TestIngressRoute(cluster, true, false, "traefik.containo.us/v1alpha1")
		})
	}

	It("\nUpgrade via SUC", func() {
		_ = specs.TestUpgradeClusterSUC(cluster, k8sClient, flags.SUCUpgradeVersion.String())
	})

	It("Checks Node status post-upgrade", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			assert.NodeAssertVersionUpgraded(),
		)
	})

	It("Checks Pod status post-upgrade", func() {
		specs.TestPodStatus(
			cluster,
			nil,
			assert.PodAssertReady())
	})

	It("Verifies ClusterIP Service post-upgrade", func() {
		specs.TestServiceClusterIP(false, true)
	})

	It("Verifies NodePort Service post-upgrade", func() {
		specs.TestServiceNodePort(false, true)
	})

	It("Verifies Ingress post-upgrade", func() {
		specs.TestIngress(false, true)
	})

	It("Verifies Daemonset post-upgrade", func() {
		specs.TestDaemonset(false, true)
	})

	It("Verifies DNS Access post-upgrade", func() {
		specs.TestDNSAccess(true, true)
	})

	if cluster.Config.Product == "rke2" {
		It("Verifies Snapshot Webhook after upgrade", func() {
			err := specs.TestSnapshotWebhook(true)
			Expect(err).To(HaveOccurred())
		})
	}

	if cluster.Config.Product == "k3s" {
		It("Verifies LoadBalancer Service after upgrade", func() {
			specs.TestServiceLoadBalancer(false, true)
		})

		It("Verifies Local Path Provisioner storage after upgrade", func() {
			specs.TestLocalPathProvisionerStorage(cluster, false, true)
		})

		It("Verifies Traefik IngressRoute after upgrade using old GKV", func() {
			specs.TestIngressRoute(cluster, false, true, "traefik.containo.us/v1alpha1")
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
