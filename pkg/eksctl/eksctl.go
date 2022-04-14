package eksctl

import (
	"eksdemo/pkg/aws"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
)

const minVersion = "0.93.0"

func GetClusterName(cluster string) string {
	return fmt.Sprintf("%s.%s.eksctl.io", cluster, aws.Region())
}

func TagNamePrefix(clusterName string) string {
	return fmt.Sprintf("eksctl-%s-cluster/", clusterName)
}

func CheckVersion() error {
	errmsg := fmt.Errorf("eksdemo requires eksctl version %s or later", minVersion)

	eksctlVersion, err := exec.Command("eksctl", "version").Output()
	if err != nil {
		return errmsg
	}

	v, err := version.NewVersion(strings.TrimSpace(string(eksctlVersion)))
	if err != nil {
		return fmt.Errorf("unable to parse eksctl version: %s", err)
	}

	if v.LessThan(version.Must(version.NewVersion(minVersion))) {
		return errmsg
	}

	return nil
}
