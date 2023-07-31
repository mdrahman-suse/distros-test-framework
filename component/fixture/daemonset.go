package fixture

import (
	// "github.com/rancher/distros-test-framework/cmd"
	"github.com/rancher/distros-test-framework/lib/shared"

	. "github.com/onsi/gomega"
)

func TestDaemonset(deployWorkload bool) {
	if deployWorkload {
		_, err := shared.ManageWorkload(
			"create",
			"daemonset.yaml",
			// cmd.ServiceFlag.ClusterConfig.Arch.String(),
		)
		Expect(err).NotTo(HaveOccurred(),
			"Daemonset manifest not deployed")
	}
	nodes, _ := shared.ParseNodes(false)
	pods, _ := shared.ParsePods(false)

	Eventually(func(g Gomega) int {
		return shared.CountOfStringInSlice("test-daemonset", pods)
	}, "420s", "10s").Should(Equal(len(nodes)),
		"Daemonset pod count does not match node count")
}
