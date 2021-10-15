package application

import "sort"

type DefaultVersion interface {
	LatestVersion(clusterVersion string) string
	PreviousVersion(clusterVersion string) string
	LatestString() string
	PreviousString() string
}

type LatestPrevious struct {
	Latest   string
	Previous string
}

func (v *LatestPrevious) LatestString() string {
	return v.Latest
}

func (v *LatestPrevious) PreviousString() string {
	return v.Previous
}

func (v *LatestPrevious) LatestVersion(clusterVersion string) string {
	return v.Latest
}

func (v *LatestPrevious) PreviousVersion(clusterVersion string) string {
	return v.Previous
}

type KubernetesVersionDependent struct {
	Latest   map[string]string
	Previous map[string]string
}

func (v *KubernetesVersionDependent) LatestString() string {
	keys := make([]string, 0, len(v.Latest))
	for k := range v.Latest {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	versionList := ""
	for _, k := range keys {
		if versionList != "" {
			versionList = versionList + "|"
		}
		versionList = versionList + v.Latest[k]
	}
	return versionList
}

func (v *KubernetesVersionDependent) PreviousString() string {
	keys := make([]string, 0, len(v.Previous))
	for k := range v.Previous {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	versionList := ""
	for _, k := range keys {
		if versionList != "" {
			versionList = versionList + "|"
		}
		versionList = versionList + v.Previous[k]
	}
	return versionList
}

func (v *KubernetesVersionDependent) LatestVersion(clusterVersion string) string {
	return v.Latest[clusterVersion]
}

func (v *KubernetesVersionDependent) PreviousVersion(clusterVersion string) string {
	return v.Previous[clusterVersion]
}
