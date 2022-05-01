package provider

import (
	"context"
	"errors"
	"fmt"
)

const (
	PROVIDER_GITHUB = "github"
	PROVIDER_GITLAB = "gitlab"
)

type Provider interface {
	GetName() string
	GetOrganizations(filter []string) ([]string, error)
	GetOrganizationRepositories(org string) ([]Repository, error)
}

type ProviderOptions struct {
	Context     context.Context
	EndpointUrl string
	AccessToken string
}

type Repository struct {
	Provider      Provider
	Owner         string
	Path          string
	Name          string
	CloneURL      string
	SSHUrl        string
	DefaultBranch string
}

func NewProvider(pType string, options ProviderOptions) (Provider, error) {
	if options.Context == nil {
		return nil, errors.New("Missing context.")
	}

	switch pType {
	case PROVIDER_GITHUB:
		return newGithubProviderClient(options)
	case PROVIDER_GITLAB:
		return newGitlabProviderClient(options)
	case "":
		return nil, errors.New("Provider should be specify.")
	default:
		return nil, errors.New("Unknow provider.")
	}
}

func (r Repository) Fullname() string {
	return r.Owner + "/" + r.Name
}

func appendError(errorList error, err error) error {
	if errorList == nil {
		return errors.New(err.Error())
	}

	return fmt.Errorf("%w; %s", errorList, err.Error())
}
