package main

import (
	"context"
	"flag"

	vacuum "github.com/jdecool/github-vacuum/internal"
	"github.com/jdecool/github-vacuum/internal/output"
	"github.com/jdecool/github-vacuum/internal/provider"
)

func main() {
	var (
		providerType        string
		providerEndpoint    string
		providerAccessToken string
		outputFormat        string
		outputFolder        string
		orgsFilter          = []string{}
	)

	appendOrg := func(org string) error {
		orgsFilter = append(orgsFilter, org)
		return nil
	}

	flag.StringVar(&providerType, "provider", "", "")
	flag.StringVar(&providerEndpoint, "provider-endpoint", "", "")
	flag.StringVar(&providerAccessToken, "provider-access-token", "", "")
	flag.StringVar(&outputFormat, "output", output.OUTPUT_FILESYSTEM, "")
	flag.StringVar(&outputFolder, "output-folder", "", "")
	flag.Func("org", "", appendOrg)
	flag.Parse()

	provider, err := provider.NewProvider(providerType, provider.ProviderOptions{
		Context:     context.Background(),
		EndpointUrl: providerEndpoint,
		AccessToken: providerAccessToken,
	})
	if err != nil {
		panic(err)
	}

	output, err := output.NewOutput(outputFormat, output.OutputOptions{
		Folder: outputFolder,
	})
	if err != nil {
		panic(err)
	}

	err = vacuum.Handle(provider, output, orgsFilter)
	if err != nil {
		panic(err)
	}
}
