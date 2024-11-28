package shared

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func setTerraformOptions(product, provider, module string) (*terraform.Options, string, error) {
	_, callerFilePath, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(callerFilePath), "..")

	providerDir, err := filepath.Abs(dir +
		fmt.Sprintf("/config/%s.tfvars", provider))
	if err != nil {
		return nil, "", fmt.Errorf("invalid provider: %s", provider)
	}
	LogLevel("info", "Set Provider path: %s", providerDir)

	varDir, err := filepath.Abs(dir + "/config/distro.tfvars")
	if err != nil {
		return nil, "", fmt.Errorf("invalid varDir path: %s", varDir)
	}
	LogLevel("info", "Set VarDir path: %s", varDir)

	// checking if module is empty, use the product as module
	if module == "" {
		module = product
	}

	tfDir, err := filepath.Abs(dir + "/modules/" + provider + "/" + module)
	if err != nil {
		return nil, "", fmt.Errorf("no module found in path: %s", tfDir)
	}
	LogLevel("info", "Set TFDir path: %s", tfDir)

	tfOpts := &terraform.Options{
		TerraformDir: tfDir,
		VarFiles:     []string{varDir, providerDir},
	}

	mergedDir, _ := filepath.Abs(dir + "/config/merged.tfvars")
	err = MergedFileContents([]string{providerDir, varDir}, mergedDir)
	if err != nil {
		return nil, "", fmt.Errorf("files not merged: %s", err)
	}

	return tfOpts, mergedDir, nil
}

func loadTFconfig(
	t *testing.T,
	product, module,
	provider, varDir string,
	tfOpts *terraform.Options,
) (*Cluster, error) {
	c := &Cluster{}

	LogLevel("info", "Loading provider config....")
	loadProviderConfig(t, tfOpts, c, provider, varDir)

	if module == "" {
		KubeConfigFile = terraform.Output(t, tfOpts, "kubeconfig")
		c.FQDN = terraform.Output(t, tfOpts, "Route53_info")
	}

	c.ServerIPs = strings.Split(terraform.Output(t, tfOpts, "master_ips"), ",")
	rawAgentIPs := terraform.Output(t, tfOpts, "worker_ips")
	if rawAgentIPs != "" {
		c.AgentIPs = strings.Split(rawAgentIPs, ",")
	}

	c.Config.Arch = terraform.GetVariableAsStringFromVarFile(t, varDir, "arch")
	c.Config.Product = product
	c.Config.ServerFlags = terraform.GetVariableAsStringFromVarFile(t, varDir, "server_flags")
	c.Config.DataStore = terraform.GetVariableAsStringFromVarFile(t, varDir, "datastore_type")
	if c.Config.DataStore == "external" {
		c.Config.ExternalDb = terraform.GetVariableAsStringFromVarFile(t, varDir, "external_db")
		c.Config.RenderedTemplate = terraform.Output(t, tfOpts, "rendered_template")
	}
	c.Config.Version = terraform.GetVariableAsStringFromVarFile(t, varDir, "install_version")

	return c, nil
}

func loadAWSConfig(t *testing.T, varDir string, c *Cluster) {
	c.AwsEc2.AccessKey = terraform.GetVariableAsStringFromVarFile(t, varDir, "access_key")
	c.AwsEc2.AwsUser = terraform.GetVariableAsStringFromVarFile(t, varDir, "aws_user")
	c.AwsEc2.Ami = terraform.GetVariableAsStringFromVarFile(t, varDir, "aws_ami")
	c.AwsEc2.Region = terraform.GetVariableAsStringFromVarFile(t, varDir, "region")
	c.AwsEc2.VolumeSize = terraform.GetVariableAsStringFromVarFile(t, varDir, "volume_size")
	c.AwsEc2.InstanceClass = terraform.GetVariableAsStringFromVarFile(t, varDir, "ec2_instance_class")
	c.AwsEc2.Subnets = terraform.GetVariableAsStringFromVarFile(t, varDir, "subnets")
	c.AwsEc2.AvailabilityZone = terraform.GetVariableAsStringFromVarFile(t, varDir, "availability_zone")
	c.AwsEc2.SgId = terraform.GetVariableAsStringFromVarFile(t, varDir, "sg_id")
	c.AwsEc2.KeyName = terraform.GetVariableAsStringFromVarFile(t, varDir, "key_name")
}

func loadvSphereConfig(t *testing.T, varDir string, c *Cluster) {
	panic("Not Implemented")
}

func addSplitRole(t *testing.T, varDir string, numServers int) (int, error) {
	splitRoles := terraform.GetVariableAsStringFromVarFile(t, varDir, "split_roles")
	if splitRoles == "true" {
		etcdNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(
			t,
			varDir,
			"etcd_only_nodes",
		))
		if err != nil {
			return 0, fmt.Errorf("error getting etcd_only_nodes %w", err)
		}
		etcdCpNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(
			t,
			varDir,
			"etcd_cp_nodes",
		))
		if err != nil {
			return 0, fmt.Errorf("error getting etcd_cp_nodes %w", err)
		}
		etcdWorkerNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(
			t,
			varDir,
			"etcd_worker_nodes",
		))
		if err != nil {
			return 0, fmt.Errorf("error getting etcd_worker_nodes %w", err)
		}
		cpNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(
			t,
			varDir,
			"cp_only_nodes",
		))
		if err != nil {
			return 0, fmt.Errorf("error getting cp_only_nodes %w", err)
		}
		cpWorkerNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(
			t,
			varDir,
			"cp_worker_nodes",
		))
		if err != nil {
			return 0, fmt.Errorf("error getting cp_worker_nodes %w", err)
		}
		numServers = numServers + etcdNodes + etcdCpNodes + etcdWorkerNodes + cpNodes + cpWorkerNodes
	}

	return numServers, nil
}

func loadProviderConfig(
	t *testing.T,
	tfOpts *terraform.Options,
	c *Cluster,
	provider, varDir string) {
	if provider == "aws" {
		LogLevel("info", "Loading aws config....")
		loadAWSConfig(t, varDir, c)
		if product == "rke2" {
			winAgentCount, err := terraform.GetVariableAsStringFromVarFileE(t, varDir, "no_of_windows_worker_nodes")
			if err != nil {
				numWinAgents, _ := strconv.Atoi(winAgentCount)
				if numWinAgents > 0 {
					LogLevel("info", "Loading Windows TF Config...")
					winAgentIPs := terraform.Output(t, tfOpts, "windows_worker_ips")
					c.WinAgentIPs = strings.Split(winAgentIPs, ",")
					c.NumWinAgents = numWinAgents
				}
			}
		}

		// TODO: Figure out logic to suppress for non-bastion clusters
		bastionCountVar, err := terraform.GetVariableAsStringFromVarFileE(t, varDir, "no_of_bastion_nodes")
		if err != nil {
			bastionCount, _ := strconv.Atoi(bastionCountVar)
			if bastionCount > 0 {
				LogLevel("info", "Loading bastion configs....")
				c.BastionConfig.PublicIPv4Addr = terraform.Output(t, tfOpts, "bastion_ip")
				c.BastionConfig.PublicDNS = terraform.Output(t, tfOpts, "bastion_dns")
			}
		}
	}
	if provider == "vsphere" {
		LogLevel("info", "Loading vsphere config....")
		loadvSphereConfig(t, varDir, c)
	}
}
