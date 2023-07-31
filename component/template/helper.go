package template

import (
	"fmt"
	"strings"
	"sync"

	"github.com/rancher/distros-test-framework/component/fixture"
	"github.com/rancher/distros-test-framework/lib/shared"
)

// upgradeVersion upgrades the version and update the expected value
func upgradeVersion(template VersionTestTemplate, product string, version string) error {
	err := fixture.TestUpgradeClusterManually(product, version)
	if err != nil {
		return err
	}

	for i := range template.TestCombination.Run {
		template.TestCombination.Run[i].ExpectedValue =
			template.TestCombination.Run[i].ExpectedValueUpgrade
	}

	return nil
}

// checkVersion checks the version and processes tests
func checkVersion(v VersionTestTemplate, product string) error {
	ips, err := getIPs()
	if err != nil {
		return fmt.Errorf("failed to get IPs: %v", err)
	}

	var wg sync.WaitGroup
	errorChanList := make(
		chan error,
		len(ips)*(len(v.TestCombination.Run)),
	)

	processTestCombination(errorChanList, &wg, product, ips, *v.TestCombination)

	wg.Wait()
	close(errorChanList)

	for errorChan := range errorChanList {
		if errorChan != nil {
			return errorChan
		}
	}

	if v.TestConfig != nil {
		TestCaseWrapper(v)
	}

	return nil
}

// getIPs gets the IPs of the nodes
func getIPs() (ips []string, err error) {
	ips = shared.FetchNodeExternalIP()
	return ips, nil
}

// AddTestCases returns the test case based on the name to be used as cmd.
func AddTestCases(names []string) ([]TestCase, error) {
	var testCases []TestCase

	testCase := map[string]TestCase{
		"TestDaemonset":                   fixture.TestDaemonset,
		"TestIngress":                     fixture.TestIngress,
		"TestDnsAccess":                   fixture.TestDnsAccess,
		"TestLocalPathProvisionerStorage": fixture.TestLocalPathProvisionerStorage,
		"TestServiceClusterIp":            fixture.TestServiceClusterIp,
		"TestServiceNodePort":             fixture.TestServiceNodePort,
		"TestServiceLoadBalancer":         fixture.TestServiceLoadBalancer,
		"TestCoredns":                     fixture.TestCoredns,
	}

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			testCases = append(testCases, func(deployWorkload bool) {})
		} else if test, ok := testCase[name]; ok {
			testCases = append(testCases, test)
		} else {
			return nil, fmt.Errorf("invalid test case name")
		}
	}

	return testCases, nil
}
