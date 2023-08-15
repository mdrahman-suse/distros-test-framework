package factory

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gruntwork-io/terratest/modules/terraform"
	//"github.com/rancher/distros-test-framework/config"
	"github.com/rancher/distros-test-framework/shared"

	. "github.com/onsi/ginkgo/v2"
)

var (
	once      sync.Once
	singleton *Cluster
)

type Cluster struct {
	Status           string
	ServerIPs        []string
	AgentIPs         []string
	NumServers       int
	NumAgents        int
	ProductType      string
	ArchType         string
	KubeConfigFile   string
	K3SCluster       K3SCluster
	RKE2Cluster      RKE2Cluster
}

type K3SCluster struct {
	DataStoreType     string
	ExternalDb        string
	RenderedTemplate  string
}

type RKE2Cluster struct {
	WinAgentIPs   []string
	NumWinAgents  int
}

// func loadConfig() (*config.ProductConfig, error) {
// 	cfg, err := config.LoadConfigEnv("./config")
// 	if err != nil {
// 		return nil, fmt.Errorf("error loading env config: %w", err)
// 	}

// 	return &cfg, nil
// }

func addTerraformOptions(g GinkgoTInterface) (*terraform.Options, string, error) {
	var varDir string
	var tfDir string
	var err error

	moduleDir := shared.BasePath() + "/distros-test-framework/modules/"
	varDir, err = filepath.Abs(moduleDir + "local.tfvars")
	if err != nil {
		return nil, "", err
	}

	ProductName := terraform.GetVariableAsStringFromVarFile(g, varDir, "product_name")
	if ProductName != "k3s" && ProductName != "rke2" {
		return nil, "", fmt.Errorf("invalid product: ", ProductName)
	}

	tfDir, err = filepath.Abs(moduleDir + ProductName)
	if err != nil {
		return nil, "", err
	}
	
	terraformOptions := &terraform.Options{
		TerraformDir: tfDir,
		VarFiles:     []string{varDir},
		EnvVars:      map[string]string{"product_name": ProductName},
	}

	return terraformOptions, varDir, nil
}

func addClusterConfig(g GinkgoTInterface, varDir string,
	terraformOptions *terraform.Options,
) (*Cluster, error) {

	cluster := &Cluster{}

	NumServers, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "no_of_server_nodes"))
	if err != nil {
		return nil, err
	}

	NumAgents, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "no_of_worker_nodes"))
	if err != nil {
		return nil, err
	}

	NumServers, err = addSplitRole(g, varDir, NumServers)
	if err != nil {
		return nil, err
	}

	cluster.NumServers = NumServers
	cluster.NumAgents = NumAgents
	cluster.ProductType = terraformOptions.EnvVars["product_name"]
	shared.Product = cluster.ProductType

	if cluster.ProductType == "k3s" {
		cluster.K3SCluster.DataStoreType = terraform.GetVariableAsStringFromVarFile(g, varDir, "datastore_type")
		if cluster.K3SCluster.DataStoreType == "" {
			cluster.K3SCluster.ExternalDb = terraform.GetVariableAsStringFromVarFile(g, varDir, "external_db")
			cluster.K3SCluster.RenderedTemplate = terraform.Output(g, terraformOptions, "rendered_template")
		}
	}

	if cluster.ProductType == "rke2" {
		rawWinAgentIPs := terraform.Output(g, terraformOptions, "windows_worker_ips")
		if rawWinAgentIPs != "" {
			cluster.RKE2Cluster.WinAgentIPs = strings.Split(rawWinAgentIPs, ",")
		}
	}
	
	shared.AwsUser = terraform.GetVariableAsStringFromVarFile(g, varDir, "aws_user")
	shared.AccessKey = terraform.GetVariableAsStringFromVarFile(g, varDir, "access_key")
	shared.Arch = terraform.GetVariableAsStringFromVarFile(g, varDir, "arch")
	shared.KubeConfigFile = terraform.Output(g, terraformOptions, "kubeconfig")
	cluster.ArchType = shared.Arch
	cluster.KubeConfigFile = shared.KubeConfigFile
	
	cluster.ServerIPs = strings.Split(terraform.Output(g, terraformOptions, "master_ips"), ",")

	rawAgentIPs := terraform.Output(g, terraformOptions, "worker_ips")
	if rawAgentIPs != "" {
		cluster.AgentIPs = strings.Split(rawAgentIPs, ",")
	}
	
	return cluster, nil
}

func addSplitRole(g GinkgoTInterface, varDir string, NumServers int) (int, error) {
	splitRoles := terraform.GetVariableAsStringFromVarFile(g, varDir, "split_roles")
	if splitRoles == "true" {
		etcdNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "etcd_only_nodes"))
		if err != nil {
			return 0, err
		}
		etcdCpNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "etcd_cp_nodes"))
		if err != nil {
			return 0, err
		}
		etcdWorkerNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "etcd_worker_nodes"))
		if err != nil {
			return 0, err
		}
		cpNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "cp_only_nodes"))
		if err != nil {
			return 0, err
		}
		cpWorkerNodes, err := strconv.Atoi(terraform.GetVariableAsStringFromVarFile(g, varDir, "cp_worker_nodes"))
		if err != nil {
			return 0, err
		}
		NumServers = NumServers + etcdNodes + etcdCpNodes + etcdWorkerNodes + cpNodes + cpWorkerNodes
	}

	return NumServers, nil
}
