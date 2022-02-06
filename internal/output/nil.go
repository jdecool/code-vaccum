package output

import "github.com/jdecool/github-vacuum/internal/provider"

type nilOutputFormatter struct {
}

func newNilOutput() (*nilOutputFormatter, error) {
	return &nilOutputFormatter{}, nil
}

func (o nilOutputFormatter) Handle(r provider.Repository) {
}

func (o nilOutputFormatter) Flush() error {
	return nil
}
