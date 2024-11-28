//go:build versionbump

package versionbump

import (
	"fmt"

	"github.com/rancher/distros-test-framework/pkg/assert"
	. "github.com/rancher/distros-test-framework/pkg/customflag"
	"github.com/rancher/distros-test-framework/pkg/specs"
	. "github.com/rancher/distros-test-framework/pkg/template"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Version Bump Template Upgrade:", func() {
	It("Start Up with no issues", func() {
		specs.TestBuildCluster(cluster)
	})

	It("Validate Nodes", func() {
		specs.TestNodeStatus(
			cluster,
			assert.NodeAssertReadyStatus(),
			nil)
	})

	It("Validate Pods", func() {
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
						Cmd:                  TestMap.Cmd,
						ExpectedValue:        TestMap.ExpectedValue,
						ExpectedValueUpgrade: TestMap.ExpectedValueUpgrade,
					},
				},
			},
			InstallMode: ServiceFlag.InstallMode.String(),
			TestConfig: &TestConfig{
				TestFunc:       ConvertToTestCase(ServiceFlag.TestTemplateConfig.TestFuncs),
				ApplyWorkload:  ServiceFlag.TestTemplateConfig.ApplyWorkload,
				DeleteWorkload: ServiceFlag.TestTemplateConfig.DeleteWorkload,
				WorkloadName:   ServiceFlag.TestTemplateConfig.WorkloadName,
			},
			Description: ServiceFlag.TestTemplateConfig.Description,
		})
	})
})

var _ = AfterEach(func() {
	if CurrentSpecReport().Failed() {
		fmt.Printf("\nFAILED! %s\n\n", CurrentSpecReport().FullText())
	} else {
		fmt.Printf("\nPASSED! %s\n\n", CurrentSpecReport().FullText())
	}
})
