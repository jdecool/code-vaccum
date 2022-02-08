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

func (p githubProvider) GetName() string {
	return PROVIDER_GITHUB
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
	var r []Repository
	var errorList error

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := p.client.Repositories.ListByOrg(p.ctx, org, nil)
		if err != nil {
			errorList = appendError(errorList, err)
			continue
		}

		for _, repo := range repos {
			r = append(r, Repository{
				Provider:      p,
				Owner:         *repo.Owner.Login,
				OwnerUrl:      *repo.Owner.HTMLURL,
				Name:          *repo.Name,
				CloneURL:      *repo.CloneURL,
				DefaultBranch: *repo.DefaultBranch,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return r, errorList
}

func (p githubProvider) getAllOrganizations() ([]string, error) {
	var errorList error
	r := []string{}

	opt := &github.OrganizationsListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		orgs, resp, err := p.client.Organizations.ListAll(p.ctx, opt)
		if err != nil {
			errorList = appendError(errorList, err)
			continue
		}

		for _, o := range orgs {
			r = append(r, *o.Login)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage

	}

	return r, nil
}
