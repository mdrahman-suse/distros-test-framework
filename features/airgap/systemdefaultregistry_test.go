//go:build systemdefaultregistry

package airgap

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"
	"github.com/rancher/distros-test-framework/shared"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test Airgap Cluster with System Default Registry:", Ordered, func() {
	It("Creates bastion and private nodes", func() {
		specs.TestBuildAirgapCluster(cluster)
	})

	It("Installs and validates product on private nodes:", func() {
		specs.TestSystemDefaultRegistry(cluster, flags)
	})

	It("Validates Nodes", func() {
		specs.TestAirgapClusterNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil,
		)
	})

	It("Validates Pods", func() {
		specs.TestAirgapClusterPodStatus(
			cluster,
			assert.PodAssertRestart(),
			assert.PodAssertReady())
	})

	AfterAll(func() {
		shared.DisplayAirgapClusterDetails(cluster)
	})

	// TODO: Validate deployment, eg: cluster-ip

})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n", CurrentSpecReport().FullText())
	}
})
