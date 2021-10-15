package helm

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

type InstallConfiguration struct {
	AppVersion    string
	ChartName     string
	Namespace     string
	ReleaseName   string
	RepositoryURL string
	ValuesFile    string
}

func initialize(kubeContext, namespace string) (*action.Configuration, error) {
	// Hack to work around https://github.com/helm/helm/issues/7430
	_ = os.Setenv("HELM_KUBECONTEXT", kubeContext)
	_ = os.Setenv("HELM_NAMESPACE", namespace)
	settings := cli.New()

	// Initialize the action configuration
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", log.Printf); err != nil {
		fmt.Print("Failed Initialize the Action Config\n")
		log.Fatal(err)
	}
	return actionConfig, nil
}

func downloadChart(ic *InstallConfiguration) (*chart.Chart, error) {
	getters := getter.All(&cli.EnvSettings{})

	// Find Chart
	chartPath, _ := repo.FindChartInRepoURL(ic.RepositoryURL, ic.ChartName, "", "", "", "", getters)

	dl := downloader.ChartDownloader{
		Out:     os.Stdout,
		Getters: getters,
	}

	u, err := dl.ResolveChartVersion(chartPath, "")
	if err != nil {
		fmt.Print("Error in ResolveChartVersion\n")
		log.Fatal(err)
	}
	fmt.Printf("Installing Chart: %s\n", u)

	g, err := dl.Getters.ByScheme(u.Scheme)
	if err != nil {
		fmt.Print("Error in ByScheme\n")
		log.Fatal(err)
	}

	// Download chart archive into memory
	data, err := g.Get(chartPath, dl.Options...)
	if err != nil {
		fmt.Print("Error in Get\n")
		log.Fatal(err)
	}

	// Decompress the archive
	files, err := loader.LoadArchiveFiles(data)
	if err != nil {
		fmt.Print("Error in LoadArchiveFiles\n")
		log.Fatal(err)
	}

	// Load the chart
	chart, err := loader.LoadFiles(files)
	if err != nil {
		fmt.Print("Error in LoadFiles\n")
		log.Fatal(err)
	}
	return chart, nil
}

func Install(ic *InstallConfiguration, kubeContext string) error {
	chart, err := downloadChart(ic)
	if err != nil {
		return err
	}

	// Parse the values file
	values := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(ic.ValuesFile), &values); err != nil {
		fmt.Print("Failed to parse values file\n")
		log.Fatal(err)
	}

	actionConfig, err := initialize(kubeContext, ic.Namespace)
	if err != nil {
		return err
	}

	// Configure the install options
	instAction := action.NewInstall(actionConfig)
	instAction.Namespace = ic.Namespace
	instAction.ReleaseName = ic.ReleaseName
	instAction.CreateNamespace = true
	instAction.IsUpgrade = true
	chart.Metadata.AppVersion = ic.AppVersion

	// Install the chart
	rel, err := instAction.Run(chart, values)
	if err != nil {
		return fmt.Errorf("helm failed to install the chart: %s", err)
	}

	log.Printf("Installed Chart: %s in namespace: %s\n", rel.Name, rel.Namespace)
	return nil
}

func List(kubeContext string) ([]*release.Release, error) {
	actionConfig, err := initialize(kubeContext, "")
	if err != nil {
		return nil, err
	}

	client := action.NewList(actionConfig)
	client.AllNamespaces = true

	releases, err := client.Run()
	if (err) != nil {
		return nil, err
	}

	return releases, nil
}

func Status(kubeContext, releaseName, namespace string) (string, error) {
	actionConfig, err := initialize(kubeContext, namespace)
	if err != nil {
		return "", err
	}
	status := action.NewStatus(actionConfig)

	rel, err := status.Run(releaseName)
	if (err) != nil {
		return "", err
	}

	// strip chart metadata from the output
	rel.Chart = nil

	return "", nil
}

func Uninstall(kubeContext, releaseName, namespace string) error {
	actionConfig, err := initialize(kubeContext, namespace)
	if err != nil {
		return err
	}
	uninstall := action.NewUninstall(actionConfig)

	// Uninstall the chart
	response, err := uninstall.Run(releaseName)
	if err != nil {
		log.Fatal("Error uninstalling Helm chart: ", errors.Cause(err))
	}
	_ = response

	return nil
}
