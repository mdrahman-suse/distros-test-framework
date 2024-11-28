package template

import (
	"fmt"
	"strings"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/pkg/k8s"
	"github.com/rancher/distros-test-framework/pkg/specs"
	"github.com/rancher/distros-test-framework/shared"
)

// upgradeVersion upgrades the product version.
func upgradeVersion(template TestTemplate, k8sClient *k8s.Client, version string) error {
	cluster := shared.ClusterConfig()
	err := specs.TestUpgradeClusterManual(cluster, k8sClient, version)
	if err != nil {
		return err
	}

	updateExpectedValue(template)

	return nil
}

// updateExpectedValue updates the expected values getting the values from flag ExpectedValueUpgrade.
func updateExpectedValue(template TestTemplate) {
	for i := range template.TestCombination.Run {
		template.TestCombination.Run[i].ExpectedValue = template.TestCombination.Run[i].ExpectedValueUpgrade
	}
}

// executeTestCombination get a template and pass it to `processTestCombination`.
//
// to execute test combination on group of IPs.
func executeTestCombination(template TestTemplate) error {
	currentVersion, err := currentProductVersion()
	if err != nil {
		return shared.ReturnLogError("failed to get current version: %w", err)
	}

	ips := shared.FetchNodeExternalIPs()
	processErr := processTestCombination(ips, currentVersion, &template)
	if processErr != nil {
		return shared.ReturnLogError("failed to process test combination: %w", processErr)
	}

	if template.TestConfig != nil {
		testCaseWrapper(template)
	}

	return nil
}

// AddTestCases returns the test case based on the name to be used as customflag.
func AddTestCases(cluster *shared.Cluster, k8sClient *k8s.Client, names []string) ([]testCase, error) {
	tcs := addTestCaseMap(cluster, k8sClient)
	return processTestCaseNames(tcs, names)
}

// addTestCaseMap initializes and returns the map of test cases.
//
//nolint:revive // we want to keep the argument for visibility.
func addTestCaseMap(cluster *shared.Cluster, k8sClient *k8s.Client) map[string]testCase {
	return map[string]testCase{
		"TestDaemonset":        specs.TestDaemonset,
		"TestIngress":          specs.TestIngress,
		"TestDNSAccess":        specs.TestDNSAccess,
		"TestServiceClusterIP": specs.TestServiceClusterIP,
		"TestServiceNodePort":  specs.TestServiceNodePort,
		"TestLocalPathProvisionerStorage": func(applyWorkload, deleteWorkload bool) {
			specs.TestLocalPathProvisionerStorage(cluster, applyWorkload, deleteWorkload)
		},
		"TestServiceLoadBalancer": specs.TestServiceLoadBalancer,
		"TestInternodeConnectivityMixedOS": func(applyWorkload, deleteWorkload bool) {
			specs.TestInternodeConnectivityMixedOS(cluster, applyWorkload, deleteWorkload)
		},
		"TestSonobuoyMixedOS": func(applyWorkload, deleteWorkload bool) {
			specs.TestSonobuoyMixedOS(deleteWorkload)
		},
		"TestSelinux": func(applyWorkload, deleteWorkload bool) {
			specs.TestSelinux(cluster)
		},
		"TestSelinuxSpcT": func(applyWorkload, deleteWorkload bool) {
			specs.TestSelinuxSpcT(cluster)
		},
		"TestUninstallPolicy": func(applyWorkload, deleteWorkload bool) {
			specs.TestUninstallPolicy(cluster)
		},
		"TestSelinuxContext": func(applyWorkload, deleteWorkload bool) {
			specs.TestSelinuxContext(cluster)
		},
		"TestIngressRoute": func(applyWorkload, deleteWorkload bool) {
			specs.TestIngressRoute(cluster, applyWorkload, deleteWorkload, "traefik.io/v1alpha1")
		},
		"TestCertRotate": func(applyWorkload, deleteWorkload bool) {
			specs.TestCertRotate(cluster)
		},
		"TestSecretsEncryption": func(applyWorkload, deleteWorkload bool) {
			specs.TestSecretsEncryption()
		},
		"TestRestartService": func(applyWorkload, deleteWorkload bool) {
			specs.TestRestartService(cluster)
		},
		"TestClusterReset": func(applyWorkload, deleteWorkload bool) {
			specs.TestClusterReset(cluster, k8sClient)
		},
	}
}

// processTestCaseNames processes the test case names and returns the corresponding test cases.
//
//nolint:revive // we want to keep the argument for visibility.
func processTestCaseNames(tcs map[string]testCase, names []string) ([]testCase, error) {
	var testCases []testCase

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			testCases = append(testCases, func(applyWorkload, deleteWorkload bool) {})
		} else if test, ok := tcs[name]; ok {
			testCases = append(testCases, test)
		} else {
			return nil, shared.ReturnLogError("invalid test case name")
		}
	}

	return testCases, nil
}

func currentProductVersion() (string, error) {
	_, version, err := shared.Product()
	if err != nil {
		return "", shared.ReturnLogError("failed to get product: %w", err)
	}
	shared.LogLevel("info", "\n\n%v", version)

	return version, nil
}

func ComponentsBumpResults() {
	product, version, err := shared.Product()
	if err != nil {
		return
	}

	var components []string
	for _, result := range assert.Results {
		if product == "rke2" {
			components = []string{"flannel", "calico", "ingressController", "coredns", "metricsServer", "etcd",
				"containerd", "runc"}
		} else {
			components = []string{"flannel", "coredns", "metricsServer", "etcd", "cniPlugins", "traefik", "local-path",
				"containerd", "klipper", "runc"}
		}
		for _, component := range components {
			if strings.Contains(result.Command, component) {
				fmt.Printf("\n---------------------\nResults from %s on version: %s\n``` \n%v\n ```\n---------------------"+
					"\n\n\n", component, version, result)
			}
		}
		fmt.Printf("\n---------------------\nResults from %s\n``` \n%v\n ```\n---------------------\n\n\n",
			result.Command, result)
	}
}
