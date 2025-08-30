package output

import (
	"errors"

	"github.com/jdecool/github-vacuum/internal/provider"
)

const (
	OUTPUT_FILESYSTEM = "filesystem"
	OUTPUT_NIL        = "nil"
	OUTPUT_REPO       = "repo"
)

type Output interface {
	Handle(r provider.Repository)
	Flush() error
}

type OutputOptions struct {
	Folder     string
	SSHKeyPath string
}

func NewOutput(format string, options OutputOptions) (Output, error) {
	switch format {
	case OUTPUT_FILESYSTEM:
		return newFilesystemOutput(FilesystemOptions{
			Folder:     options.Folder,
			SSHKeyPath: options.SSHKeyPath,
		})
	case OUTPUT_NIL:
		return newNilOutput()
	case OUTPUT_REPO:
		return newRepoOutput()
	default:
		return nil, errors.New("Unknown output format.")
	}
}
