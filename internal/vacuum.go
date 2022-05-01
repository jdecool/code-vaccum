package vacuum

import (
	"errors"
	"fmt"

	"github.com/jdecool/github-vacuum/internal/output"
	"github.com/jdecool/github-vacuum/internal/provider"
	log "github.com/sirupsen/logrus"
)

func Handle(p provider.Provider, o output.Output, orgsFilter []string) error {
	var errorList error

	orgs, err := p.GetOrganizations(orgsFilter)
	log.Debugf("Found %d organization(s)", len(orgs))

	if err != nil {
		log.Debug(err.Error())
		return err
	}

	for _, org := range orgs {
		log.Debugf("Processing organization: %s", org)

		repos, err := p.GetOrganizationRepositories(org)
		log.Debugf("Found %d repository(ies)", len(repos))

		if err != nil {
			log.Debug(err.Error())
			errorList = appendError(errorList, err)
		}

		for _, repo := range repos {
			log.Debugf("Processing repository: %s", repo.Fullname())
			o.Handle(repo)
		}
	}

	if err = o.Flush(); err != nil {
		log.Debug(err.Error())
		errorList = appendError(errorList, err)
	}

	return errorList
}

func appendError(errorList error, err error) error {
	if errorList == nil {
		return errors.New(err.Error())
	}

	return fmt.Errorf("%w; %s", errorList, err.Error())
}
