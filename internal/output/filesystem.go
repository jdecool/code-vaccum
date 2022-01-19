package output

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/jdecool/github-vacuum/internal/provider"
)

type filesystemOutputFormatter struct {
	opts FilesystemOptions
}

type FilesystemOptions struct {
	Folder string
}

func newFilesystemOutput(opts FilesystemOptions) (*filesystemOutputFormatter, error) {
	return &filesystemOutputFormatter{opts}, nil
}

func (o filesystemOutputFormatter) Handle(r provider.Repository) {
	path := o.opts.Folder
	if strings.TrimSpace(path) != "" {
		path += "/"
	}
	path += r.Owner + "/" + r.Name

	git.PlainClone(path, false, &git.CloneOptions{
		URL: r.CloneURL,
	})
}
