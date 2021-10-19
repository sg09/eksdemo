package eksctl

import (
	"eksdemo/pkg/resource"
	"eksdemo/pkg/template"
	"encoding/json"
	"fmt"
	"strings"
)

const EksctlHeader = `
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: {{ .ClusterName }}
  region: {{ .Region }}
  tags: 
    eksdemo: v0.1.0-alpha
{{- if .KubernetesVersion }}
  version: {{ .KubernetesVersion | printf "%q" }}
{{ end }}
`

type ResourceManager struct {
	Resource      string
	Template      template.Template
	ApproveCreate bool
	ApproveDelete bool
	*IamAuth
}

type IamAuth struct {
	Arn      string
	Groups   []string
	Username string
}

func (e *ResourceManager) Create(options resource.Options) error {
	switch e.Resource {
	case "iamidentitymapping":
		return e.CreateIamAuth(options)
	default:
		return e.CreateWithConfigFile(options)
	}
}

func (e *ResourceManager) CreateIamAuth(options resource.Options) error {
	t := template.TextTemplate{Template: e.IamAuth.Arn}
	renderedArn, err := t.Render(options)

	if err != nil {
		return fmt.Errorf("failed to render ARN %q: %s", e.IamAuth.Arn, err)
	}

	if exists, err := e.IamAuthExists(renderedArn, options.Common().ClusterName); err != nil {
		return err
	} else if exists {
		fmt.Printf("Iam Auth %q already exists\n", options.Common().Name)
		return nil
	}

	args := []string{
		"create",
		e.Resource,
		"--arn",
		renderedArn,
		"--username",
		e.IamAuth.Username,
		"--group",
		strings.Join(e.IamAuth.Groups, ","),
		"--cluster",
		options.Common().ClusterName,
	}

	return Command(args, "")
}

func (e *ResourceManager) CreateWithConfigFile(options resource.Options) error {
	eksctlConfig, err := e.Template.Render(options)

	if err != nil {
		return err
	}

	args := []string{
		"create",
		e.Resource,
		"-f",
		"-",
	}

	if e.ApproveCreate {
		args = append(args, "--approve")
	}

	return Command(args, eksctlConfig)
}

func (e *ResourceManager) Delete(options resource.Options) error {
	switch e.Resource {
	case "iamidentitymapping":
		return e.DeleteIamAuth(options)
	default:
		return e.DeleteWithConfigFile(options)
	}
}

func (e *ResourceManager) DeleteIamAuth(options resource.Options) error {
	t := template.TextTemplate{Template: e.IamAuth.Arn}
	renderedArn, err := t.Render(options)

	if err != nil {
		return fmt.Errorf("failed to render ARN %q: %s", e.IamAuth.Arn, err)
	}

	args := []string{
		"delete",
		e.Resource,
		"--arn",
		renderedArn,
		"--cluster",
		options.Common().ClusterName,
	}

	return Command(args, "")
}

func (e *ResourceManager) DeleteWithConfigFile(options resource.Options) error {
	options.PrepForDelete()
	eksctlConfig, err := e.Template.Render(options)

	if err != nil {
		return err
	}

	args := []string{
		"delete",
		e.Resource,
		"-f",
		"-",
	}

	if e.ApproveDelete {
		args = append(args, "--approve")
	}

	return Command(args, eksctlConfig)
}

func (e *ResourceManager) IamAuthExists(renderedArn, cluster string) (bool, error) {
	args := []string{
		"get",
		e.Resource,
		"--arn",
		renderedArn,
		"--cluster",
		cluster,
		"-o",
		"json",
	}

	result, err := CommandWithResult(args, "")

	// eksctl writes to stderr if no results found
	if err != nil {
		return false, nil
	}

	var jsonObjs interface{}
	json.Unmarshal([]byte(result), &jsonObjs)
	jsonSlice, ok := jsonObjs.([]interface{})
	if !ok {
		return false, fmt.Errorf("failed to parse eksctl json output")
	}

	return len(jsonSlice) > 0, nil
}
