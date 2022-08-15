package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func ECRDescribeRepositories(name string) ([]*ecr.Repository, error) {
	sess := GetSession()
	svc := ecr.New(sess)

	repositories := []*ecr.Repository{}
	input := &ecr.DescribeRepositoriesInput{}
	pageNum := 0

	if name != "" {
		input.RepositoryNames = aws.StringSlice([]string{name})
	}

	err := svc.DescribeRepositoriesPages(input,
		func(page *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
			pageNum++
			repositories = append(repositories, page.Repositories...)
			return pageNum <= maxPages
		},
	)

	return repositories, err
}
