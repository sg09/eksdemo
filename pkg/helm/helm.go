package helm

import (
	"fmt"
	"log"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

type Helm struct {
	AppVersion    string
	ChartName     string
	Namespace     string
	PostRenderer  postrender.PostRenderer
	ReleaseName   string
	RepositoryURL string
	ValuesFile    string
	Wait          bool
}

func initialize(kubeContext, namespace string) (*action.Configuration, error) {
	// Hack to work around https://github.com/helm/helm/issues/7430
	_ = os.Setenv("HELM_KUBECONTEXT", kubeContext)
	_ = os.Setenv("HELM_NAMESPACE", namespace)
	settings := cli.New()

	// Initialize the action configuration
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", log.Printf); err != nil {
		return nil, fmt.Errorf("failed to initialize helm action config: %w", err)
	}
	return actionConfig, nil
}

func (h *Helm) DownloadChart() (*chart.Chart, error) {
	getters := getter.All(&cli.EnvSettings{})

	// Find Chart
	chartPath, err := repo.FindChartInRepoURL(h.RepositoryURL, h.ChartName, "", "", "", "", getters)
	if err != nil {
		return nil, err
	}

	dl := downloader.ChartDownloader{
		Out:     os.Stdout,
		Getters: getters,
	}

	u, err := dl.ResolveChartVersion(chartPath, "")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Downloading Chart: %s\n", u)

	g, err := dl.Getters.ByScheme(u.Scheme)
	if err != nil {
		return nil, err
	}

	// Download chart archive into memory
	data, err := g.Get(chartPath, dl.Options...)
	if err != nil {
		return nil, err
	}

	// Decompress the archive
	files, err := loader.LoadArchiveFiles(data)
	if err != nil {
		return nil, err
	}

	// Load the chart
	chart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

func (h *Helm) Install(chart *chart.Chart, kubeContext string) error {
	// Parse the values file
	values := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(h.ValuesFile), &values); err != nil {
		return fmt.Errorf("failed to parse values file: %w", err)
	}

	actionConfig, err := initialize(kubeContext, h.Namespace)
	if err != nil {
		return err
	}

	// Configure the install options
	instAction := action.NewInstall(actionConfig)
	instAction.Namespace = h.Namespace
	instAction.ReleaseName = h.ReleaseName
	instAction.CreateNamespace = true
	instAction.IsUpgrade = true
	instAction.PostRenderer = h.PostRenderer
	instAction.Wait = h.Wait
	instAction.Timeout = 300 * time.Second
	chart.Metadata.AppVersion = h.AppVersion

	// Install the chart
	fmt.Println("Helm installing...")
	rel, err := instAction.Run(chart, values)
	if err != nil {
		return fmt.Errorf("helm install failed: %s", err)
	}

	fmt.Printf("Installed: %s in namespace: %s\n", rel.Name, rel.Namespace)
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
	_, err = uninstall.Run(releaseName)
	if err != nil {
		return fmt.Errorf("failed uninstalling chart: %w", err)
	}

	return nil
}
