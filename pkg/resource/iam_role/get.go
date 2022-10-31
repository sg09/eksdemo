package iam_role

import (
	"eksdemo/pkg/aws"
	"eksdemo/pkg/printer"
	"eksdemo/pkg/resource"
	"eksdemo/pkg/resource/iam_oidc"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type Getter struct {
	iamClient  *aws.IAMClient
	oidcGetter *iam_oidc.Getter
}

func NewGetter(iamClient *aws.IAMClient) *Getter {
	return &Getter{iamClient, iam_oidc.NewGetter(iamClient)}
}

func (g *Getter) Init() {
	if g.iamClient == nil {
		g.iamClient = aws.NewIAMClient()
	}
	g.oidcGetter = iam_oidc.NewGetter(g.iamClient)
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	roleOptions, ok := options.(*IamRoleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to IamRoleOptions")
	}

	var role *types.Role
	var roles []types.Role
	var err error

	if name != "" {
		role, err = g.GetRoleByName(name)
	} else if roleOptions.Cluster != nil {
		roles, err = g.GetIrsaRolesForCluster(roleOptions.Cluster, roleOptions.LastUsed)
	} else {
		roles, err = g.GetAllRoles(roleOptions.All, roleOptions.LastUsed)
	}

	if err != nil {
		return err
	}

	if role != nil {
		roles = []types.Role{*role}
	}

	return output.Print(os.Stdout, NewPrinter(roles, roleOptions.LastUsed))
}

func (g *Getter) GetAllRoles(includeServiceRoles, getRoleDetails bool) (roles []types.Role, err error) {
	roles, err = g.iamClient.ListRoles()
	if err != nil {
		return nil, err
	}

	if !includeServiceRoles {
		filtered := []types.Role{}
		for _, r := range roles {
			if !strings.HasPrefix(awssdk.ToString(r.Path), "/aws-service-role/") {
				filtered = append(filtered, r)
			}
		}
		roles = filtered
	}

	if getRoleDetails {
		return g.getDetailedRoles(roles)
	}

	return roles, nil
}

func (g *Getter) GetIrsaRolesForCluster(cluster *ekstypes.Cluster, getRoleDetails bool) ([]types.Role, error) {
	oidc, err := g.oidcGetter.GetOidcProviderByCluster(cluster)
	if err != nil {
		return []types.Role{}, err
	}

	roles, err := g.GetAllRoles(false, getRoleDetails)
	if err != nil {
		return []types.Role{}, err
	}

	irsaRoles := []types.Role{}
	providerUrlEscaped := url.QueryEscape(awssdk.ToString(oidc.Url))

	fmt.Println(providerUrlEscaped)

	for _, r := range roles {
		if strings.Contains(awssdk.ToString(r.AssumeRolePolicyDocument), providerUrlEscaped) {
			irsaRoles = append(irsaRoles, r)
		}
	}

	return irsaRoles, nil
}

func (g *Getter) GetRoleByName(name string) (*types.Role, error) {
	role, err := g.iamClient.GetRole(name)

	if err != nil {
		var nsee *types.NoSuchEntityException
		if errors.As(err, &nsee) {
			return nil, resource.NotFoundError(fmt.Sprintf("iam-role %q not found", name))
		}
		return nil, err
	}

	return role, nil
}

func (g *Getter) getDetailedRoles(roles []types.Role) ([]types.Role, error) {
	detailedRoles := make([]types.Role, 0, len(roles))
	for _, r := range roles {
		role, err := g.iamClient.GetRole(awssdk.ToString(r.RoleName))
		if err != nil {
			return []types.Role{}, err
		}

		detailedRoles = append(detailedRoles, *role)
	}

	return detailedRoles, nil
}
