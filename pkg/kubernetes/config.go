package kubernetes

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// ClientConfig is used to make it easy to get an api server client
//
// type ClientConfig interface {
// 	// RawConfig returns the merged result of all overrides
// 	RawConfig() (clientcmdapi.Config, error)
// 	// ClientConfig returns a complete client config
// 	ClientConfig() (*restclient.Config, error)
// 	// Namespace returns the namespace resulting from the merged
// 	// result of all overrides and a boolean indicating if it was
// 	// overridden
// 	Namespace() (string, bool, error)
// 	// ConfigAccess returns the rules for loading/persisting the config.
// 	ConfigAccess() ConfigAccess
// }

// Raw is clientcmdapi.Config -- represents kubeconfig
func GetCurrentContextClusterURL() string {
	raw, err := GetKubeconfig()
	if err != nil {
		return ""
	}

	if err := clientcmdapi.MinifyConfig(raw); err != nil {
		return ""
	}

	return raw.Clusters[raw.Contexts[raw.CurrentContext].Cluster].Server
}

func GetKubeContextForCluster(cluster *eks.Cluster) (string, error) {
	raw, err := GetKubeconfig()
	if err != nil {
		return "", err
	}

	found := ""

	for name, context := range raw.Contexts {
		if _, ok := raw.Clusters[context.Cluster]; ok {
			if raw.Clusters[context.Cluster].Server == aws.StringValue(cluster.Endpoint) {
				found = name
				break
			}
		}
	}

	return found, nil
}

func GetKubeconfig() (*clientcmdapi.Config, error) {
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	raw, err := config.RawConfig()
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

// TODO: Refactor below -- used only by pkg/cluster/list.go
// TODO: use-context command should use this library

func GetClientConfig() clientcmdapi.Config {
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	raw, err := config.RawConfig()
	if err != nil {
		log.Fatal(err)
	}

	return raw
}

func GetClientConfigCurrentCluster() string {
	raw := GetClientConfig()

	if len(raw.Contexts) == 0 {
		return ""
	}
	return raw.Contexts[raw.CurrentContext].Cluster
}
