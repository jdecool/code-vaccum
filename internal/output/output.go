package output

import (
	"errors"

	"github.com/jdecool/github-vacuum/internal/provider"
)

const (
	OUTPUT_FILESYSTEM = "filesystem"
	OUTPUT_NIL        = "nil"
)

type Output interface {
	Handle(r provider.Repository)
	Flush()
}

type OutputOptions struct {
	Folder string
}

func NewOutput(format string, options OutputOptions) (Output, error) {
	switch format {
	case OUTPUT_FILESYSTEM:
		return newFilesystemOutput(FilesystemOptions{
			Folder: options.Folder,
		})
	case OUTPUT_NIL:
		return newNilOutput()
	default:
		return nil, errors.New("Unknown output format.")
	}
}
