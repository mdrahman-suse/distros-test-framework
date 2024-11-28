//go:build flannel

package versionbump

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"

	"github.com/rancher/distros-test-framework/pkg/assert"
	. "github.com/rancher/distros-test-framework/pkg/customflag"
	"github.com/rancher/distros-test-framework/pkg/specs"
	. "github.com/rancher/distros-test-framework/pkg/template"
)

var flannelCmd string

var _ = Describe("Flannel Version bump:", func() {
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

	It("Test flannel version bump", func() {
		flannelCmd = "kubectl get node -o yaml : | grep 'hardened-flannel' -A1"
		if cluster.Config.Product == "k3s" {
			flannelCmd = "/var/lib/rancher/k3s/data/current/bin/flannel"
		}

		Template(TestTemplate{
			TestCombination: &RunCmd{
				Run: []TestMapConfig{
					{
						Cmd:                  flannelCmd,
						ExpectedValue:        TestMap.ExpectedValue,
						ExpectedValueUpgrade: TestMap.ExpectedValueUpgrade,
					},
				},
			},
			InstallMode: ServiceFlag.InstallMode.String(),
		})
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
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
