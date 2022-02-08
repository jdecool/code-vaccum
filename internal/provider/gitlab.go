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

func (p gitlabProvider) GetName() string {
	return PROVIDER_GITLAB
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
	var r []Repository
	var errorList error

	opt := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := p.client.Groups.ListGroupProjects(org, nil)
		if err != nil {
			errorList = appendError(errorList, err)
			continue
		}

		for _, repo := range repos {
			r = append(r, Repository{
				Provider:      p,
				Owner:         org,
				OwnerUrl:      repo.WebURL[0:strings.LastIndex(repo.WebURL, "/")],
				Name:          repo.Name,
				CloneURL:      repo.HTTPURLToRepo,
				DefaultBranch: repo.DefaultBranch,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return r, nil
}

func (p gitlabProvider) getAllOrganizations() ([]string, error) {
	var errorList error
	r := []string{}

	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}

	for {
		groups, resp, err := p.client.Groups.ListGroups(nil)
		if err != nil {
			errorList = appendError(errorList, err)
			continue
		}

		for _, g := range groups {
			r = append(r, g.Name)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return r, errorList
}
