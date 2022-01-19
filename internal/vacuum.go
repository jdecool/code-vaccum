package vacuum

import (
	"errors"
	"fmt"

	"github.com/jdecool/github-vacuum/internal/output"
	"github.com/jdecool/github-vacuum/internal/provider"
)

func Handle(p provider.Provider, o output.Output, orgsFilter []string) error {
	var errorList error

	orgs, err := p.GetOrganizations(orgsFilter)
	if err != nil {
		return err
	}

	for _, org := range orgs {
		repos, err := p.GetOrganizationRepositories(org)
		if err != nil {
			appendError(errorList, err)
		}

		for _, repo := range repos {
			fmt.Println(repo)
			o.Handle(repo)
		}
	}

	return errorList
}

func appendError(errorList error, err error) error {
	if errorList == nil {
		return errors.New(err.Error())
	}

	return fmt.Errorf("%w; %s", errorList, err.Error())
}
