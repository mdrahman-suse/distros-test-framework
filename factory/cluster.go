package factory

import (
	"fmt"
	//"path/filepath"
	//"strconv"

	"github.com/gruntwork-io/terratest/modules/terraform"
	//"github.com/rancher/distros-test-framework/shared"

	. "github.com/onsi/ginkgo/v2"
)

// NewCluster creates a new cluster and returns his values from terraform config and vars
func NewCluster(g GinkgoTInterface) (*Cluster, error) {
	terraformOptions, varDir, err := addTerraformOptions(g)
	if err != nil {
		return nil, err
	}

	fmt.Println("Creating Cluster")
	terraform.InitAndApply(g, terraformOptions)

	cluster, err := addClusterConfig(g, varDir, terraformOptions)
	if err != nil {
		return nil, err
	}

	cluster.Status = "cluster created"

	return cluster, nil
}

// GetCluster returns a singleton cluster
func GetCluster(g GinkgoTInterface) *Cluster {
	var err error
	once.Do(func() {
		singleton, err = NewCluster(g)
		if err != nil {
			g.Errorf("error getting cluster: %v", err)
		}
	})
	return singleton
}

// DestroyCluster destroys the cluster and returns a message
func DestroyCluster(g GinkgoTInterface) (string, error) {

	// cfg, err := loadConfig()
	// if err != nil {
	// 	return "", fmt.Errorf("error loading config: %w", err)
	// }
	terraformOptions, _ , err := addTerraformOptions(g)
	if err != nil {
		return "", err
	}

	terraform.Destroy(g, terraformOptions)

	return "cluster destroyed", nil
}
