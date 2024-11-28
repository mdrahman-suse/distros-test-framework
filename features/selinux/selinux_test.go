package selinux

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/customflag"
	"github.com/rancher/distros-test-framework/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
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

	It("Validate selinux is enabled", func() {
		specs.TestSelinuxEnabled(cluster)
	})

	It("Validate container, server and selinux version", func() {
		specs.TestSelinux(cluster)
	})

	It("Validate container security", func() {
		specs.TestSelinuxSpcT(cluster)
	})

	It("Validate context", func() {
		specs.TestSelinuxContext(cluster)
	})

	if customflag.ServiceFlag.InstallMode.String() != "" {
		It("Upgrade manual", func() {
			_ = specs.TestUpgradeClusterManual(cluster, k8sClient, customflag.ServiceFlag.InstallMode.String())
		})

		It("Validate Nodes Post upgrade", func() {
			specs.TestNodeStatus(
				cluster,
				assert.NodeAssertReadyStatus(),
				assert.NodeAssertVersionTypeUpgrade(&customflag.ServiceFlag),
			)
		})

		It("Validate Pods Post upgrade", func() {
			specs.TestPodStatus(
				cluster,
				assert.PodAssertRestart(),
				assert.PodAssertReady())
		})

		It("Validate selinux is enabled Post upgrade", func() {
			specs.TestSelinuxEnabled(cluster)
		})

		It("Validate container, server and selinux version Post upgrade", func() {
			specs.TestSelinux(cluster)
		})

		It("Validate container security Post upgrade", func() {
			specs.TestSelinuxSpcT(cluster)
		})

		It("Validate context", func() {
			specs.TestSelinuxContext(cluster)
		})
	}

	It("Validate uninstall selinux policies", func() {
		specs.TestUninstallPolicy(cluster)
	})

})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
