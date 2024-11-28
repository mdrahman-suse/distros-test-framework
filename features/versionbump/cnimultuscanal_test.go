//go:build multus

package versionbump

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	. "github.com/rancher/distros-test-framework/pkg/customflag"
	"github.com/rancher/distros-test-framework/pkg/specs"
	. "github.com/rancher/distros-test-framework/pkg/template"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Multus + canal Version bump:", func() {
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

	It("Test Bump version", func() {
		Template(TestTemplate{
			TestCombination: &RunCmd{
				Run: []TestMapConfig{
					{
						Cmd: "kubectl get node -o yaml : | grep multus-cni -A1, " +
							"kubectl -n kube-system get pods -l k8s-app=canal -o jsonpath=\"{..image}\" : " +
							"| awk '{for(i=1;i<=NF;i++) if($i ~ /calico/) print $i}', " +
							" kubectl -n kube-system get pods -l k8s-app=canal -o jsonpath=\"{..image}\" : " +
							"| awk '{for(i=1;i<=NF;i++) if($i ~ /flannel/) print $i}' , " +
							"kubectl get pods -n kube-system : | grep multus | awk '{print $1} {print $3}'",
						ExpectedValue:        TestMap.ExpectedValue,
						ExpectedValueUpgrade: TestMap.ExpectedValueUpgrade,
					},
				},
			},
			InstallMode: ServiceFlag.InstallMode.String(),
		})
	})

	It("Verifies dns access", func() {
		specs.TestDNSAccess(true, true)
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
