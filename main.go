package main

import (
	"context"
	"flag"

	vacuum "github.com/jdecool/github-vacuum/internal"
	"github.com/jdecool/github-vacuum/internal/output"
	"github.com/jdecool/github-vacuum/internal/provider"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		providerType        string
		providerEndpoint    string
		providerAccessToken string
		outputFormat        string
		outputFolder        string
		sshKeyPath          string
		orgsFilter          = []string{}
		usernamesFilter     = []string{}
		debug               bool
		quiet               bool
	)

	appendOrg := func(org string) error {
		orgsFilter = append(orgsFilter, org)
		return nil
	}

	appendUsername := func(username string) error {
		usernamesFilter = append(usernamesFilter, username)
		return nil
	}

	flag.StringVar(&providerType, "provider", "", "")
	flag.StringVar(&providerEndpoint, "provider-endpoint", "", "")
	flag.StringVar(&providerAccessToken, "provider-access-token", "", "")
	flag.StringVar(&outputFormat, "output", output.OUTPUT_FILESYSTEM, "")
	flag.StringVar(&outputFolder, "output-folder", "", "")
	flag.StringVar(&sshKeyPath, "ssh-key", "", "")
	flag.BoolVar(&debug, "debug", false, "")
	flag.BoolVar(&quiet, "quiet", false, "")
	flag.Func("org", "", appendOrg)
	flag.Func("username", "", appendUsername)
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	} else if quiet {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	p, err := provider.NewProvider(providerType, provider.ProviderOptions{
		Context:     context.Background(),
		EndpointUrl: providerEndpoint,
		AccessToken: providerAccessToken,
	})
	if err != nil {
		panic(err)
	}

	o, err := output.NewOutput(outputFormat, output.OutputOptions{
		Folder:     outputFolder,
		SSHKeyPath: sshKeyPath,
	})
	if err != nil {
		panic(err)
	}

	err = vacuum.Handle(p, o, orgsFilter, usernamesFilter)
	if err != nil {
		panic(err)
	}
}
