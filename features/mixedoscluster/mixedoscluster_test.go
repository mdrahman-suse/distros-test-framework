package mixedoscluster

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Test: Mixed OS Cluster", func() {

	It("Starts Up with no issues", func() {
		specs.TestBuildCluster(cluster)
	})

	It("Validates Node", func() {
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

	It("Validates internode connectivity over the vxlan tunnel", func() {
		specs.TestInternodeConnectivityMixedOS(cluster, true, true)
	})

	It("Validates cluster by running sonobuoy mixed OS plugin", func() {
		specs.TestSonobuoyMixedOS(true)
	})
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
