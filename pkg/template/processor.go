package template

import (
	"strings"

	"github.com/rancher/distros-test-framework/pkg/assert"
	"github.com/rancher/distros-test-framework/shared"
)

// processTestCombination processes the test combination on a group of IPs,sending values to processCmds.
func processTestCombination(
	ips []string,
	currentVersion string,
	t *TestTemplate,
) error {
	if t.TestCombination.Run != nil {
		for _, testMap := range t.TestCombination.Run {
			cmds := strings.Split(testMap.Cmd, ",")
			expectedValues := strings.Split(testMap.ExpectedValue, ",")

			if strings.Contains(testMap.Cmd, "etcd ") {
				nodes, err := shared.GetNodesByRoles("etcd")
				if err != nil {
					shared.LogLevel("error", "error from getting nodes by roles: %w\n", err)
					return err
				}

				var externalIPs []string
				var ip string
				for _, n := range nodes {
					ip = n.ExternalIP
					externalIPs = append(externalIPs, ip)
				}

				ips = externalIPs
			}

			for _, ip := range ips {
				if processErr := processCmds(ip, cmds, expectedValues, currentVersion); processErr != nil {
					return shared.ReturnLogError("error from processCmds: %w", processErr)
				}
			}
		}
	}

	return nil
}

// processCmds runs the tests per ips using processOnNode and processOnHost validation.
func processCmds(
	ip string,
	cmds []string,
	expectedValues []string,
	currentProductVersion string,
) error {
	// range over the cmds only cause expectedValues arrives here on the same length.
	for i, c := range cmds {
		expectedValue := strings.TrimSpace(strings.Trim(expectedValues[i], "\""))
		cmd := strings.TrimSpace(strings.Trim(c, "\""))

		if strings.Contains(c, "kubectl") || strings.HasPrefix(cmd, "helm") {
			processHostErr := processOnHost(cmd, expectedValue, currentProductVersion)
			if processHostErr != nil {
				return shared.ReturnLogError("error from processOnHost: %w", processHostErr)
			}
		} else {
			processNodeErr := processOnNode(cmd, expectedValue, ip, currentProductVersion)
			if processNodeErr != nil {
				return shared.ReturnLogError("error from processOnNode: %w", processNodeErr)
			}
		}
	}

	return nil
}

// processOnNode runs the test on the node calling ValidateOnNode.
func processOnNode(cmd, expectedValue, ip, currentProductVersion string) error {
	if currentProductVersion == "" {
		shared.LogLevel("error", "error getting current version, is empty\n")
		return shared.ReturnLogError("error getting current version, is empty\n")
	}

	shared.LogLevel("debug", "Version Check: %s\nIP Address: %s\nCommand to Execute: "+
		"%s\nExecution Location: Node\nExpected Value: %s\n",
		currentProductVersion, ip, cmd, expectedValue)

	err := assert.ValidateOnNode(ip, cmd, expectedValue)
	if err != nil {
		return shared.ReturnLogError("error from validate on node: %w\n", err)
	}

	return nil
}

// processOnHost runs the test on the host calling ValidateOnHost.
func processOnHost(cmd, expectedValue, currentProductVersion string) error {
	if currentProductVersion == "" {
		return shared.ReturnLogError("error getting current version, is empty\n")
	}

	shared.LogLevel("debug", "Version Check: %s\nCommand to Execute: %s\nExecution Location: Host\nExpected Value: %s\n",
		currentProductVersion, cmd, expectedValue)

	kubeconfigFlag := " --kubeconfig=" + shared.KubeConfigFile
	var fullCmd string
	if strings.Contains(cmd, ":") {
		fullCmd = shared.JoinCommands(cmd, kubeconfigFlag)
	} else {
		fullCmd = cmd + kubeconfigFlag
	}

	fullCmd = strings.ReplaceAll(fullCmd, `"`, "")
	err := assert.ValidateOnHost(fullCmd, expectedValue)
	if err != nil {
		return shared.ReturnLogError("error from validate on host: %w\n", err)
	}

	return nil
}
