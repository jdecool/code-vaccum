package output

import (
	"github.com/jdecool/github-vacuum/internal/provider"
	log "github.com/sirupsen/logrus"
)

type nilOutputFormatter struct {
}

func newNilOutput() (*nilOutputFormatter, error) {
	return &nilOutputFormatter{}, nil
}

func (o nilOutputFormatter) Handle(r provider.Repository) {
	log.Infof("Processing repository: %s", r.Fullname())
}

func (o nilOutputFormatter) Flush() error {
	return nil
}
