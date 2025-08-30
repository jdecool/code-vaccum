package vacuum

import (
	"errors"
	"fmt"
	"time"

	"github.com/jdecool/github-vacuum/internal/output"
	"github.com/jdecool/github-vacuum/internal/provider"
	log "github.com/sirupsen/logrus"
)

func Handle(p provider.Provider, o output.Output, orgsFilter []string, usernamesFilter []string) error {
	var errorList error
	startTime := time.Now()
	totalRepos := 0
	processedRepos := 0

	log.Infof("Starting vacuum operation with provider: %s", p.GetName())

	if len(orgsFilter) > 0 || len(usernamesFilter) == 0 {
		log.Infof("Processing %d organization(s)...", len(orgsFilter))

		orgs, err := p.GetOrganizations(orgsFilter)
		log.Infof("Found %d organization(s)", len(orgs))

		if err != nil {
			log.Error("Error fetching organizations: ", err.Error())
			errorList = appendError(errorList, err)
		}

		for orgIdx, org := range orgs {
			log.Infof("[%d/%d] Processing organization: %s", orgIdx+1, len(orgs), org)

			repos, err := p.GetOrganizationRepositories(org)
			repoCount := len(repos)
			totalRepos += repoCount
			log.Infof("Found %d repository(ies) in organization %s", repoCount, org)

			if err != nil {
				log.Error("Error fetching repositories for org ", org, ": ", err.Error())
				errorList = appendError(errorList, err)
			}

			for repoIdx, repo := range repos {
				processedRepos++
				log.Infof("[%d/%d] Processing repository: %s (%d/%d total)", repoIdx+1, repoCount, repo.Fullname(), processedRepos, totalRepos)
				o.Handle(repo)
			}
		}
	}

	for userIdx, username := range usernamesFilter {
		log.Infof("[%d/%d] Processing user: %s", userIdx+1, len(usernamesFilter), username)

		repos, err := p.GetUserRepositories(username)
		repoCount := len(repos)
		totalRepos += repoCount
		log.Infof("Found %d repository(ies) for user %s", repoCount, username)

		if err != nil {
			log.Error("Error fetching repositories for user ", username, ": ", err.Error())
			errorList = appendError(errorList, err)
		}

		for repoIdx, repo := range repos {
			processedRepos++
			log.Infof("[%d/%d] Processing repository: %s (%d/%d total)", repoIdx+1, repoCount, repo.Fullname(), processedRepos, totalRepos)
			o.Handle(repo)
		}
	}

	log.Info("Flushing output...")
	if err := o.Flush(); err != nil {
		log.Error("Error flushing output: ", err.Error())
		errorList = appendError(errorList, err)
	}

	duration := time.Since(startTime)
	log.Infof("Vacuum operation completed in %v", duration)
	log.Infof("Processed %d repositories total", processedRepos)

	if errorList != nil {
		log.Warn("Operation completed with errors")
	} else {
		log.Info("Operation completed successfully")
	}

	return errorList
}

func appendError(errorList error, err error) error {
	if errorList == nil {
		return errors.New(err.Error())
	}

	return fmt.Errorf("%w; %s", errorList, err.Error())
}
