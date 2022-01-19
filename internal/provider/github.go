package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubProvider struct {
	ctx    context.Context
	client *github.Client
}

func newGithubProviderClient(options ProviderOptions) (*githubProvider, error) {
	c, err := createGithubClient(options)
	if err != nil {
		return nil, err
	}

	return &githubProvider{
		options.Context,
		c,
	}, nil
}

func createGithubClient(options ProviderOptions) (*github.Client, error) {
	httpClient := createHttpClient(options)

	githubUrl := options.EndpointUrl
	if strings.TrimSpace(githubUrl) == "" {
		return github.NewClient(httpClient), nil
	}

	c, err := github.NewEnterpriseClient(githubUrl, githubUrl, httpClient)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func createHttpClient(options ProviderOptions) *http.Client {
	if strings.TrimSpace(options.AccessToken) == "" {
		return nil
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: options.AccessToken,
		},
	)

	return oauth2.NewClient(options.Context, tokenSource)
}

func (p githubProvider) GetOrganizations(filter []string) ([]string, error) {
	if len(filter) == 0 {
		return p.getAllOrganizations()
	}

	var errorList error

	r := []string{}
	for _, org := range filter {
		o, _, err := p.client.Organizations.Get(p.ctx, org)
		if err != nil {
			errorList = appendError(errorList, err)
			continue
		}

		r = append(r, *o.Name)
	}

	return r, errorList
}

func (p githubProvider) GetOrganizationRepositories(org string) ([]Repository, error) {
	repos, _, err := p.client.Repositories.ListByOrg(p.ctx, org, nil)
	if err != nil {
		return []Repository{}, err
	}

	var r []Repository
	for _, repo := range repos {
		r = append(r, Repository{
			Owner:         *repo.Owner.Login,
			Name:          *repo.Name,
			CloneURL:      *repo.CloneURL,
			DefaultBranch: *repo.DefaultBranch,
		})
	}

	return r, nil
}

func (p githubProvider) getAllOrganizations() ([]string, error) {
	orgs, _, err := p.client.Organizations.ListAll(p.ctx, nil)
	if err != nil {
		return []string{}, err
	}

	r := []string{}
	for _, o := range orgs {
		r = append(r, *o.Login)
	}

	return r, nil
}
