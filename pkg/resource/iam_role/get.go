package iam_role

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

type Getter struct{}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	roleOptions, ok := options.(*IamRoleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to IamRoleOptions")
	}

	roles, err := g.GetRoles(name)
	if err != nil {
		return err
	}

	// If listing roles without the --all flag, filter out Service Roles
	if name == "" && !roleOptions.All {
		filtered := []*iam.Role{}
		for _, r := range roles {
			if !strings.HasPrefix(aws.StringValue(r.Path), "/aws-service-role/") {
				filtered = append(filtered, r)
			}
		}
		roles = filtered
	}

	// If getting a role by name, or using the --last-used flag, call GetRole
	// because ListRoles doesn't include the last used date (or include tags)
	if name != "" || roleOptions.LastUsed {
		detailedRoles := make([]*iam.Role, 0, len(roles))

		for _, r := range roles {
			role, err := aws.IamGetRole(aws.StringValue(r.RoleName))
			if err != nil {
				return err
			}

			detailedRoles = append(detailedRoles, role)
		}
		roles = detailedRoles
	}

	return output.Print(os.Stdout, NewPrinter(roles, roleOptions.LastUsed))
}

func (g *Getter) GetRoles(name string) (roles []*iam.Role, err error) {
	roles, err = aws.IamListRoles()
	if err != nil {
		return nil, err
	}

	if name != "" {
		filtered := []*iam.Role{}
		for _, r := range roles {
			if strings.EqualFold(aws.StringValue(r.RoleName), name) {
				filtered = append(filtered, r)
			}
		}
		roles = filtered
	}

	return roles, nil
}
