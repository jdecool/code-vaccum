package output

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/jdecool/github-vacuum/internal/provider"
)

type repoOuputFormatter struct {
	processedRemotes map[string]manifestRemote
	data             manifest
}

type manifest struct {
	Remotes  []manifestRemote  `xml:"remote"`
	Projects []manifestProject `xml:"project"`
}

type manifestRemote struct {
	Name  string `xml:"name,attr"`
	Fetch string `xml:"fetch,attr"`
}

type manifestProject struct {
	Name     string `xml:"name,attr"`
	Remote   string `xml:"remote,attr"`
	Path     string `xml:"path,attr"`
	Revision string `xml:"revision,attr"`
}

func newRepoOutput() (*repoOuputFormatter, error) {
	return &repoOuputFormatter{
		processedRemotes: map[string]manifestRemote{},
		data:             manifest{},
	}, nil
}

func (o *repoOuputFormatter) Handle(repo provider.Repository) {
	p := repo.Provider

	remote, exists := o.processedRemotes[repo.OwnerUrl]
	if !exists {
		remote = o.data.AddRemote(p, repo)
		o.processedRemotes[repo.OwnerUrl] = remote
	}

	o.data.AddProject(remote, repo)
}

func (o repoOuputFormatter) Flush() error {
	xml, err := xml.MarshalIndent(o.data, " ", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", xml)

	return nil
}

func (m *manifest) AddRemote(p provider.Provider, r provider.Repository) manifestRemote {
	remote := manifestRemote{
		Name:  strings.ToLower(p.GetName()) + "-" + strings.ToLower(r.Owner),
		Fetch: r.OwnerUrl,
	}

	m.Remotes = append(m.Remotes, remote)

	return remote
}

func (m *manifest) AddProject(remote manifestRemote, repo provider.Repository) {
	project := manifestProject{
		Name:     repo.Name,
		Remote:   remote.Name,
		Path:     repo.Owner + "/" + repo.Name,
		Revision: repo.DefaultBranch,
	}

	m.Projects = append(m.Projects, project)
}
