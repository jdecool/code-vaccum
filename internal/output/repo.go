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
	Revision string `xml:"revision,attr"`
}

func newRepoOutput() (*repoOuputFormatter, error) {
	return &repoOuputFormatter{
		processedRemotes: map[string]manifestRemote{},
		data:             manifest{},
	}, nil
}

func (o *repoOuputFormatter) Handle(repo provider.Repository) {
	remoteName := repo.Path[0:strings.Index(repo.Path, "/")]
	remoteUrl := repo.SSHUrl[0:strings.Index(repo.SSHUrl, remoteName)] + remoteName
	if !strings.HasPrefix(remoteUrl, "ssh://") {
		remoteUrl = strings.Replace(remoteUrl, ":", "/", 1)
		remoteUrl = "ssh://" + remoteUrl
	}

	remote, exists := o.processedRemotes[remoteName]
	if !exists {
		remote = o.data.AddRemote(repo.Provider, remoteName, remoteUrl)
		o.processedRemotes[remoteName] = remote
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

func (m *manifest) AddRemote(p provider.Provider, name string, url string) manifestRemote {
	remote := manifestRemote{
		Name:  name,
		Fetch: url,
	}

	m.Remotes = append(m.Remotes, remote)

	return remote
}

func (m *manifest) AddProject(remote manifestRemote, repo provider.Repository) {
	project := manifestProject{
		Name:     remote.Name + "/" + strings.Replace(repo.Path, remote.Name+"/", "", 1),
		Remote:   remote.Name,
		Revision: repo.DefaultBranch,
	}

	m.Projects = append(m.Projects, project)
}
