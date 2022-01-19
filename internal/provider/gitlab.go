package provider

import (
	"strings"

	"github.com/xanzy/go-gitlab"
)

type gitlabProvider struct {
	client *gitlab.Client
}

func newGitlabProviderClient(options ProviderOptions) (*gitlabProvider, error) {
	c, err := gitlab.NewClient(options.AccessToken)
	if err != nil {
		return nil, err
	}

	return &gitlabProvider{
		c,
	}, nil
}

func createGitlabClient(options ProviderOptions) (*gitlab.Client, error) {
	if strings.TrimSpace(options.EndpointUrl) == "" {
		return gitlab.NewClient(options.AccessToken)
	}

	return gitlab.NewClient(options.AccessToken, gitlab.WithBaseURL(options.EndpointUrl))
}

func (p gitlabProvider) GetOrganizations(filter []string) ([]string, error) {
	if len(filter) == 0 {
		return p.getAllOrganizations()
	}

	var errorList error

	r := []string{}
	for _, org := range filter {
		g, _, err := p.client.Groups.GetGroup(org, nil)
		if err != nil {
			errorList = appendError(errorList, err)
			continue
		}

		r = append(r, g.Name)
	}

	return r, errorList
}

func (p gitlabProvider) GetOrganizationRepositories(org string) ([]Repository, error) {
	repos, _, err := p.client.Groups.ListGroupProjects(org, nil)
	if err != nil {
		return []Repository{}, err
	}

	var r []Repository
	for _, repo := range repos {
		r = append(r, Repository{
			Owner:         org,
			Name:          repo.Name,
			CloneURL:      repo.HTTPURLToRepo,
			DefaultBranch: repo.DefaultBranch,
		})
	}

	return r, nil
}

func (p gitlabProvider) getAllOrganizations() ([]string, error) {
	groups, _, err := p.client.Groups.ListGroups(nil)
	if err != nil {
		return []string{}, err
	}

	r := []string{}
	for _, g := range groups {
		r = append(r, g.Name)
	}

	return r, nil
}
